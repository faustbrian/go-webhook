package webhook

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// IdempotencyHeader is the local signed application idempotency field, not a receiver guarantee.
const IdempotencyHeader = "Idempotency-Key"

var (
	// ErrDeliveryFailed categorizes delivery construction/transport/status failures.
	ErrDeliveryFailed = errors.New("webhook delivery failed")
	// ErrEndpointRejected categorizes endpoint admission refusal.
	ErrEndpointRejected = errors.New("webhook endpoint rejected")
	// ErrResponseTooLarge means the response sentinel read exceeded its byte bound.
	ErrResponseTooLarge = errors.New("webhook response too large")
	// ErrFanOutLimit means input delivery count exceeded MaxFanOut.
	ErrFanOutLimit = errors.New("webhook fan-out limit exceeded")
)

// HTTPDoer is implemented by http.Client and compatible client wrappers.
type HTTPDoer interface {
	// Do must honor request context and return a non-nil Response when err is nil.
	// The Deliverer closes response bodies, including response-plus-error cases.
	// Caller configuration must enforce redirects/proxy/egress and bounded I/O.
	Do(request *http.Request) (*http.Response, error)
}

// EndpointPolicy validates an endpoint immediately before every attempt.
type EndpointPolicy interface {
	// Validate admits or rejects a shallow per-attempt URL copy; implementations
	// must be concurrency-safe, honor ctx and return promptly. Validation alone
	// cannot prevent a different address being selected by an unsafe HTTP client.
	Validate(ctx context.Context, endpoint *url.URL) error
}

// EndpointPolicyFunc adapts a function into an endpoint policy.
type EndpointPolicyFunc func(ctx context.Context, endpoint *url.URL) error

// Validate implements EndpointPolicy.
func (f EndpointPolicyFunc) Validate(ctx context.Context, endpoint *url.URL) error {
	return f(ctx, endpoint)
}

// SleepFunc provides injectable delivery backoff. Implementations must honor
// ctx, be safe for concurrent calls and bound their work. Their errors return
// directly; callbacks are not forcibly canceled or panic-contained.
type SleepFunc func(ctx context.Context, duration time.Duration) error

// FailureClassification is a stable delivery outcome category.
type FailureClassification string

const (
	// FailureNone marks a locally successful attempt.
	FailureNone FailureClassification = "none"
	// FailureRetryable marks a failure scheduled for another attempt.
	FailureRetryable FailureClassification = "retryable"
	// FailureTerminal marks a nonretryable/construction/body/canceled failure.
	FailureTerminal FailureClassification = "terminal"
	// FailureExhausted marks a retryable failure with no attempt allowance remaining.
	FailureExhausted FailureClassification = "exhausted"
)

// RetryPolicy bounds attempts and exponential backoff.
type RetryPolicy struct {
	// MaxAttempts must be positive at construction and includes the first try; no key forces one.
	MaxAttempts int
	// BaseDelay is nonnegative initial fallback delay, zero allowed.
	BaseDelay time.Duration
	// MaxDelay caps both Retry-After and exponential delay and must be at least BaseDelay.
	MaxDelay time.Duration
}

// Delay returns a Retry-After delay when valid, otherwise exponential backoff.
// Negative or inverted delay bounds return zero. Attempts at or below one use
// the base delay; exponential work stops once the configured cap is reached.
func (p RetryPolicy) Delay(attempt int, now time.Time, retryAfter string) time.Duration {
	if p.BaseDelay < 0 || p.MaxDelay < p.BaseDelay {
		return 0
	}
	delay := time.Duration(0)
	if seconds, err := strconv.ParseInt(retryAfter, 10, 64); err == nil && seconds >= 0 {
		// Bound seconds before multiplication so even MaxInt64 Retry-After
		// values cannot overflow the signed duration domain.
		if seconds > int64(p.MaxDelay/time.Second) {
			delay = p.MaxDelay
		} else {
			delay = time.Duration(seconds) * time.Second
		}
	} else if parsed, err := http.ParseTime(retryAfter); err == nil && parsed.After(now) {
		delay = parsed.Sub(now)
	} else {
		delay = p.BaseDelay
		if attempt > 1 {
			for remaining := attempt - 1; remaining > 0 && delay > 0 && delay < p.MaxDelay; remaining-- {
				if delay > p.MaxDelay-delay {
					delay = p.MaxDelay
				} else {
					delay += delay
				}
			}
		}
	}

	return min(delay, p.MaxDelay)
}

// DeliveryAttempt records one admitted attempt, including request-construction
// failure before HTTP dispatch. Core diagnostics exclude payloads, secrets,
// signatures, endpoint query strings, or sensitive response headers.
type DeliveryAttempt struct {
	// ID is a caller-generated nonempty attempt ID; no length or privacy validation is added.
	ID string
	// Number is one-based within this delivery.
	Number int
	// StartedAt is sampled before request construction.
	StartedAt time.Time
	// CompletedAt is sampled after Do or construction failure, before response body read.
	CompletedAt time.Time
	// StatusCode is returned HTTP status, zero on construction/transport failure.
	StatusCode int
	// Classification records local retry policy, not receiver-side effect certainty.
	Classification FailureClassification
	// RetryAfter is planned bounded backoff; it does not prove sleep completed.
	RetryAfter time.Duration
	// Diagnostic is a fixed core diagnostic; caller-constructed values are not sanitized.
	Diagnostic string
}

// DeliveryResult contains per-call evidence. Attempts/response bytes are owned
// slices but may also be shared with DeadLetter callbacks. IDs and response data
// are application/endpoint-controlled and not sanitized; do not log blindly.
type DeliveryResult struct {
	// ID is a caller-generated delivery identity, possibly empty on pre-admission failure.
	ID string
	// EventID is copied caller attribution, not independently authenticated by outbound delivery.
	EventID string
	// Attempts is per-call evidence in attempt order, including request-construction failures.
	Attempts []DeliveryAttempt
	// ResponseBody is owned bounded bytes from the last successfully read response; it may be sensitive.
	ResponseBody []byte
}

// DeliveryRequest describes one endpoint delivery.
type DeliveryRequest struct {
	// Endpoint is borrowed and shallow-copied per attempt before EndpointPolicy.
	Endpoint *url.URL
	// Body is borrowed across attempts; do not mutate it during delivery or serialization.
	Body []byte
	// EventID must be nonempty attribution; no automatic event-ID HTTP header is emitted.
	EventID string
	// IdempotencyKey enables core retries when nonempty; receiver semantics remain application-owned.
	IdempotencyKey string
	// Headers is cloned per attempt; JSON Content-Type and explicit idempotency override matching values.
	Headers http.Header
	// Metadata is borrowed canonical signed data; bound size and avoid concurrent mutation.
	Metadata map[string]string
}

// DeadLetterFunc receives a terminal bounded result for durable handling by
// the application. The core does not implement storage or a queue. It runs
// synchronously on selected terminal paths, not all early/canceled failures.
// Result slices may alias returned evidence: do not mutate them. Hook errors
// join the primary error and may expose caller text; panics are not contained.
type DeadLetterFunc func(ctx context.Context, result DeliveryResult) error

// ReplayHook records an operator-requested replay before new delivery validation.
// It is optional, borrowed and synchronous; errors prevent delivery, panics
// propagate, and the application owns authorization and durable audit semantics.
type ReplayHook func(ctx context.Context, originalDeliveryID, eventID string) error

// DeliveryConfig configures a bounded deliverer.
type DeliveryConfig struct {
	// Client is a mandatory borrowed HTTPDoer; redirects/proxy/timeout policy is caller-owned.
	Client HTTPDoer
	// Signer is a mandatory constructed signer retained across calls.
	Signer *Signer
	// EndpointPolicy is mandatory and invoked before each try; use a validated transport for dial-time safety.
	EndpointPolicy EndpointPolicy
	// Retry must have positive attempts and nonnegative ordered delay bounds.
	Retry RetryPolicy
	// Clock is borrowed; nil uses time.Now. It must be concurrency-safe and return promptly.
	Clock func() time.Time
	// Sleep is borrowed; nil uses a context-aware timer. Honor context and bound callback work.
	Sleep SleepFunc
	// IDGenerator is mandatory, borrowed and must provide fresh nonempty IDs concurrently.
	IDGenerator func() (string, error)
	// MaxRequestBytes must be positive below MaxInt64; caps request body bytes, not metadata/headers.
	MaxRequestBytes int64
	// MaxResponseBytes must be positive below MaxInt64; caps retained response bytes plus one sentinel read.
	MaxResponseBytes int64
	// MaxFanOut must be positive and caps delivery count, not all concurrent calls or total memory.
	MaxFanOut int
	// HeaderLimits bounds generated signature fields and must have positive members.
	HeaderLimits HeaderLimits
	// DeadLetter is optional synchronous terminal-result handling, not durable storage itself.
	DeadLetter DeadLetterFunc
	// ReplayHook is optional synchronous operator replay audit, invoked before new delivery validation.
	ReplayHook ReplayHook
	// Observer is optional synchronous observation; core contains its panics, not its blocking.
	Observer Observer
}

// Deliverer signs and sends bounded webhook attempts. Use NewDeliverer; zero
// is unusable. Concurrent calls/FanOut require safe client/signer/policy/clock/
// ID/sleep/hooks. It owns no persistent worker, queue or client close lifecycle.
// Bounds are per call, not global throughput or total process memory guarantees.
type Deliverer struct {
	client           HTTPDoer
	signer           *Signer
	policy           EndpointPolicy
	retry            RetryPolicy
	clock            func() time.Time
	sleep            SleepFunc
	id               func() (string, error)
	maxRequestBytes  int64
	maxResponseBytes int64
	maxFanOut        int
	headerLimits     HeaderLimits
	deadLetter       DeadLetterFunc
	replayHook       ReplayHook
	observer         Observer
}

// NewDeliverer validates positive bounds, collaborators and retry delays before
// returning an immutable configuration owner. Nil interfaces are refused but
// typed nil is not detected. Collaborators are borrowed and not invoked; no
// timeout, redirect or proxy policy is installed on a supplied Client.
func NewDeliverer(config DeliveryConfig) (*Deliverer, error) {
	if err := validateDeliveryConfig(config); err != nil {
		return nil, err
	}
	clock := config.Clock
	if clock == nil {
		clock = time.Now
	}
	sleep := config.Sleep
	if sleep == nil {
		sleep = sleepContext
	}

	return &Deliverer{
		client:           config.Client,
		signer:           config.Signer,
		policy:           config.EndpointPolicy,
		retry:            config.Retry,
		clock:            clock,
		sleep:            sleep,
		id:               config.IDGenerator,
		maxRequestBytes:  config.MaxRequestBytes,
		maxResponseBytes: config.MaxResponseBytes,
		maxFanOut:        config.MaxFanOut,
		headerLimits:     config.HeaderLimits,
		deadLetter:       config.DeadLetter,
		replayHook:       config.ReplayHook,
		observer:         config.Observer,
	}, nil
}

func validateDeliveryConfig(config DeliveryConfig) error {
	switch {
	case config.Client == nil:
		return fmt.Errorf("%w: client is required", ErrInvalidConfiguration)
	case config.Signer == nil:
		return fmt.Errorf("%w: signer is required", ErrInvalidConfiguration)
	case config.EndpointPolicy == nil:
		return fmt.Errorf("%w: endpoint policy is required", ErrInvalidConfiguration)
	case config.IDGenerator == nil:
		return fmt.Errorf("%w: ID generator is required", ErrInvalidConfiguration)
	case config.MaxRequestBytes <= 0 || config.MaxRequestBytes == math.MaxInt64:
		return fmt.Errorf("%w: safe positive request limit is required", ErrInvalidConfiguration)
	case config.MaxResponseBytes <= 0 || config.MaxResponseBytes == math.MaxInt64:
		return fmt.Errorf("%w: safe positive response limit is required", ErrInvalidConfiguration)
	case config.MaxFanOut <= 0:
		return fmt.Errorf("%w: positive fan-out limit is required", ErrInvalidConfiguration)
	case config.HeaderLimits.MaxSignatures <= 0:
		return fmt.Errorf("%w: positive signature count limit is required", ErrInvalidConfiguration)
	case config.HeaderLimits.MaxBytes <= 0:
		return fmt.Errorf("%w: positive signature byte limit is required", ErrInvalidConfiguration)
	case config.Retry.MaxAttempts <= 0:
		return fmt.Errorf("%w: positive attempt limit is required", ErrInvalidConfiguration)
	case config.Retry.BaseDelay < 0:
		return fmt.Errorf("%w: retry base delay cannot be negative", ErrInvalidConfiguration)
	case config.Retry.MaxDelay < config.Retry.BaseDelay:
		return fmt.Errorf("%w: retry maximum delay cannot be below the base delay", ErrInvalidConfiguration)
	default:
		return nil
	}
}

// Deliver performs bounded attempts. Retries are disabled for requests that
// do not carry an explicit idempotency key. It uses POST/JSON, re-signs each try
// and admits policy before construction. 2xx succeeds; transport errors and
// 408/425/429/500/502/503/504 can retry. Returned result may contain partial
// evidence/last response on error; caller context and cooperative collaborators
// bound blocking. Observer panics are contained, lifecycle/client panics are not.
func (d *Deliverer) Deliver(ctx context.Context, delivery DeliveryRequest) (DeliveryResult, error) {
	result := DeliveryResult{EventID: delivery.EventID}
	if delivery.Endpoint == nil {
		return result, fmt.Errorf("%w: endpoint is required", ErrInvalidConfiguration)
	}
	if delivery.EventID == "" {
		return result, fmt.Errorf("%w: event ID is required", ErrInvalidConfiguration)
	}
	if int64(len(delivery.Body)) > d.maxRequestBytes {
		return result, fmt.Errorf("%w: request body exceeds limit", ErrInvalidConfiguration)
	}
	deliveryID, err := d.id()
	if err != nil {
		return result, fmt.Errorf("%w: delivery ID generation failed", ErrDeliveryFailed)
	}
	if deliveryID == "" {
		return result, fmt.Errorf("%w: delivery ID generation failed", ErrDeliveryFailed)
	}
	result.ID = deliveryID
	maxAttempts := d.retry.MaxAttempts
	if delivery.IdempotencyKey == "" {
		maxAttempts = 1
	}

	for attemptIndex := range maxAttempts {
		number := attemptIndex + 1
		if err := ctx.Err(); err != nil {
			return result, err
		}
		endpoint := *delivery.Endpoint
		if err := d.policy.Validate(ctx, &endpoint); err != nil {
			failure := fmt.Errorf("%w: policy denied attempt", ErrEndpointRejected)
			return result, d.terminalError(ctx, result, failure)
		}
		attemptID, err := d.id()
		if err != nil {
			failure := fmt.Errorf("%w: attempt ID generation failed", ErrDeliveryFailed)
			return result, d.terminalError(ctx, result, failure)
		}
		if attemptID == "" {
			failure := fmt.Errorf("%w: attempt ID generation failed", ErrDeliveryFailed)
			return result, d.terminalError(ctx, result, failure)
		}
		attempt := DeliveryAttempt{ID: attemptID, Number: number, StartedAt: d.clock().UTC(), Classification: FailureNone}
		request, err := d.request(ctx, &endpoint, delivery)
		if err != nil {
			attempt.CompletedAt = d.clock().UTC()
			attempt.Classification = FailureTerminal
			attempt.Diagnostic = "request construction failed"
			d.recordAttempt(ctx, &result, attempt)
			return result, d.terminalError(ctx, result, err)
		}

		response, doErr := d.client.Do(request)
		attempt.CompletedAt = d.clock().UTC()
		if doErr != nil {
			if response != nil && response.Body != nil {
				_ = response.Body.Close()
			}
			attempt.Diagnostic = "HTTP transport failed"
			if err := ctx.Err(); err != nil {
				attempt.Classification = FailureTerminal
				d.recordAttempt(ctx, &result, attempt)
				return result, err
			}
			if number == maxAttempts {
				attempt.Classification = FailureExhausted
				d.recordAttempt(ctx, &result, attempt)
				failure := fmt.Errorf("%w: transport attempts exhausted", ErrDeliveryFailed)
				return result, d.terminalError(ctx, result, failure)
			}
			attempt.Classification = FailureRetryable
			attempt.RetryAfter = d.retry.Delay(number, d.clock().UTC(), "")
			d.recordAttempt(ctx, &result, attempt)
			if err := d.sleep(ctx, attempt.RetryAfter); err != nil {
				return result, err
			}
			continue
		}

		body, bodyErr := readResponse(response, d.maxResponseBytes)
		attempt.StatusCode = response.StatusCode
		if bodyErr != nil {
			attempt.Classification = FailureTerminal
			attempt.Diagnostic = "response body exceeded limit or could not be read"
			d.recordAttempt(ctx, &result, attempt)
			return result, d.terminalError(ctx, result, bodyErr)
		}
		result.ResponseBody = body
		if response.StatusCode >= 200 && response.StatusCode < 300 {
			d.recordAttempt(ctx, &result, attempt)
			return result, nil
		}

		if retryableStatus(response.StatusCode) {
			if number < maxAttempts {
				attempt.Classification = FailureRetryable
				attempt.RetryAfter = d.retry.Delay(number, d.clock().UTC(), response.Header.Get("Retry-After"))
				d.recordAttempt(ctx, &result, attempt)
				if err := d.sleep(ctx, attempt.RetryAfter); err != nil {
					return result, err
				}
				continue
			}
			attempt.Classification = FailureExhausted
		} else {
			attempt.Classification = FailureTerminal
		}
		attempt.Diagnostic = "endpoint returned non-success status"
		d.recordAttempt(ctx, &result, attempt)
		failure := fmt.Errorf("%w: endpoint returned status %d", ErrDeliveryFailed, response.StatusCode)

		return result, d.terminalError(ctx, result, failure)
	}

	return result, fmt.Errorf("%w: attempt limit exhausted without a result", ErrDeliveryFailed)
}

// DeliverOnce performs at most one HTTP attempt even when the Deliverer retry
// policy allows more. Queue and outbox consumers use this to avoid nested
// retry multiplication.
func (d *Deliverer) DeliverOnce(ctx context.Context, delivery DeliveryRequest) (DeliveryResult, error) {
	single := *d
	single.retry.MaxAttempts = 1

	return single.Deliver(ctx, delivery)
}

func (d *Deliverer) request(ctx context.Context, endpoint *url.URL, delivery DeliveryRequest) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(delivery.Body))
	if err != nil {
		return nil, fmt.Errorf("%w: request creation failed", ErrDeliveryFailed)
	}
	request.Header = delivery.Headers.Clone()
	if request.Header == nil {
		request.Header = make(http.Header)
	}
	request.Header.Set("Content-Type", "application/json")
	if delivery.IdempotencyKey != "" {
		request.Header.Set(IdempotencyHeader, delivery.IdempotencyKey)
	}
	_, _, err = d.signer.SignRequest(request, RequestOptions{
		MaxBodyBytes: d.maxRequestBytes,
		HeaderLimits: d.headerLimits,
		Metadata:     delivery.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: request signing failed", ErrDeliveryFailed)
	}

	return request, nil
}

func readResponse(response *http.Response, maxBytes int64) ([]byte, error) {
	if response == nil || response.Body == nil {
		return nil, fmt.Errorf("%w: missing response body", ErrDeliveryFailed)
	}
	body, readErr := io.ReadAll(io.LimitReader(response.Body, maxBytes+1))
	closeErr := response.Body.Close()
	if readErr != nil {
		return nil, fmt.Errorf("%w: response read failed", ErrDeliveryFailed)
	}
	if closeErr != nil {
		return nil, fmt.Errorf("%w: response read failed", ErrDeliveryFailed)
	}
	if int64(len(body)) > maxBytes {
		return nil, ErrResponseTooLarge
	}

	return body, nil
}

func retryableStatus(status int) bool {
	switch status {
	case http.StatusRequestTimeout, http.StatusTooEarly, http.StatusTooManyRequests,
		http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func (d *Deliverer) recordAttempt(ctx context.Context, result *DeliveryResult, attempt DeliveryAttempt) {
	result.Attempts = append(result.Attempts, attempt)
	outcome := OutcomeFailure
	reason := ReasonTransport
	switch attempt.Classification {
	case FailureNone:
		outcome = OutcomeSuccess
		reason = ReasonNone
	case FailureRetryable:
		outcome = OutcomeRetry
		if attempt.StatusCode != 0 {
			reason = ReasonStatus
		}
	case FailureTerminal, FailureExhausted:
		if attempt.StatusCode != 0 {
			reason = ReasonStatus
		}
	}
	observeSafely(d.observer, ctx, Observation{
		Operation:      OperationDeliveryAttempt,
		Outcome:        outcome,
		Reason:         reason,
		Duration:       elapsed(func() time.Time { return attempt.CompletedAt }, attempt.StartedAt),
		StatusCode:     attempt.StatusCode,
		Attempt:        attempt.Number,
		Classification: attempt.Classification,
	})
}

func (d *Deliverer) terminalError(ctx context.Context, result DeliveryResult, primary error) error {
	if d.deadLetter == nil {
		return primary
	}
	if err := d.deadLetter(ctx, result); err != nil {
		return errors.Join(primary, fmt.Errorf("dead-letter hook failed: %w", err))
	}

	return primary
}

// FanOutResult preserves the input order of a bounded fan-out operation.
type FanOutResult struct {
	// Result is delivery evidence at the matching input index, even on per-delivery failure.
	Result DeliveryResult
	// Err is the matching delivery error; FanOut returns it here, not as aggregate error.
	Err error
}

// FanOut runs deliveries with min(workers, input count) joined workers and no
// durable queue. workers must be positive and count <= MaxFanOut. It preserves
// input order, not dispatch/completion order; delivery errors live in each result.
// It waits for all workers, including canceled calls. A collaborator ignoring
// context can delay return; no worker is detached to satisfy a timeout. Inputs
// stay borrowed for the entire call and must not be concurrently mutated.
func (d *Deliverer) FanOut(ctx context.Context, deliveries []DeliveryRequest, workers int) ([]FanOutResult, error) {
	if workers < 1 {
		return nil, fmt.Errorf("%w: positive fan-out workers required", ErrInvalidConfiguration)
	}
	if len(deliveries) > d.maxFanOut {
		return nil, ErrFanOutLimit
	}
	results := make([]FanOutResult, len(deliveries))
	if len(deliveries) == 0 {
		return results, nil
	}
	workers = min(max(workers, 1), len(deliveries))
	jobs := make(chan int, len(deliveries))
	for index := range deliveries {
		jobs <- index
	}
	close(jobs)
	var group sync.WaitGroup
	group.Add(workers)
	for range workers {
		go func() {
			defer group.Done()
			for index := range jobs {
				results[index].Result, results[index].Err = d.Deliver(ctx, deliveries[index])
			}
		}()
	}
	group.Wait()

	return results, nil
}

// Replay audits operator intent before starting a new independently identified
// delivery. It does not mutate or reuse the original attempt record. Hook nil
// permits replay without audit; the core does not look up original state or
// prove operator permission. Fresh identity depends on caller IDGenerator.
func (d *Deliverer) Replay(
	ctx context.Context,
	originalDeliveryID string,
	delivery DeliveryRequest,
) (DeliveryResult, error) {
	if originalDeliveryID == "" {
		return DeliveryResult{}, fmt.Errorf("%w: original delivery ID required", ErrInvalidConfiguration)
	}
	if d.replayHook != nil {
		if err := d.replayHook(ctx, originalDeliveryID, delivery.EventID); err != nil {
			return DeliveryResult{}, fmt.Errorf("replay hook failed: %w", err)
		}
	}

	return d.Deliver(ctx, delivery)
}

func sleepContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
