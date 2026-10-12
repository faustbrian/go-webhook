package webhook

import (
	"context"
	"errors"
	"time"
)

// Operation is a low-cardinality observed operation.
type Operation string

const (
	// OperationVerification identifies bounded HTTP verification, not plain Verify.
	OperationVerification Operation = "verification"
	// OperationReplay identifies atomic replay-store admission.
	OperationReplay Operation = "replay"
	// OperationDeliveryAttempt identifies recorded outbound attempt evidence.
	OperationDeliveryAttempt Operation = "delivery_attempt"
)

// Outcome is a low-cardinality observed result.
type Outcome string

const (
	// OutcomeSuccess means the observed operation completed without returned failure.
	OutcomeSuccess Outcome = "success"
	// OutcomeRejected means inbound verification or replay admission refused.
	OutcomeRejected Outcome = "rejected"
	// OutcomeRetry means another outbound attempt is planned.
	OutcomeRetry Outcome = "retry"
	// OutcomeFailure means a recorded outbound attempt failed.
	OutcomeFailure Outcome = "failure"
)

// Reason names the fixed core failure categories. It is string-backed, so a
// direct caller can construct arbitrary data; adapters do not normalize it.
type Reason string

const (
	// ReasonNone means no observed failure category.
	ReasonNone Reason = "none"
	// ReasonSignature categorizes malformed or unauthenticated signature data.
	ReasonSignature Reason = "signature"
	// ReasonTimestamp categorizes timestamp refusal.
	ReasonTimestamp Reason = "timestamp"
	// ReasonReplay categorizes retained replay identity.
	ReasonReplay Reason = "replay"
	// ReasonStore categorizes replay-store unavailability.
	ReasonStore Reason = "store"
	// ReasonLimit categorizes body/header/response size rejection.
	ReasonLimit Reason = "limit"
	// ReasonCanceled categorizes context cancellation/deadline errors.
	ReasonCanceled Reason = "canceled"
	// ReasonTransport categorizes delivery/transport failure without status attribution.
	ReasonTransport Reason = "transport"
	// ReasonStatus categorizes non-success endpoint status.
	ReasonStatus Reason = "status"
	// ReasonPolicy categorizes endpoint admission refusal.
	ReasonPolicy Reason = "policy"
	// ReasonInternal categorizes remaining failures without exposing error text.
	ReasonInternal Reason = "internal"
)

// Observation intentionally has no payload, signature, key, event ID,
// endpoint, header, query, or arbitrary error field.
type Observation struct {
	// Operation is a fixed core operation; direct callers can construct arbitrary string values.
	Operation Operation
	// Outcome is a fixed core result category, not durable-effect evidence.
	Outcome Outcome
	// Reason is a fixed core failure category rather than raw error text.
	Reason Reason
	// Duration is measured/clamped nonnegative by core; direct callers control it.
	Duration time.Duration
	// Algorithm is the configured algorithm for verification/replay, empty on core delivery events.
	Algorithm Algorithm
	// StatusCode is attempt HTTP status or zero; adapters reduce it to status class.
	StatusCode int
	// Attempt is the one-based delivery attempt, zero for verification/replay.
	Attempt int
	// Classification is delivery retry classification, empty for verification/replay.
	Classification FailureClassification
}

// Observer receives fixed-schema core observations synchronously. Borrowed
// implementations must support concurrent calls, bound retained state and return
// promptly. Core invocation contains panics but adds no queue/retry or timeout;
// direct calls do not inherit that protection. Arbitrary caller-created values
// are not automatically privacy-safe.
type Observer interface {
	// Observe receives one copied value; it cannot veto the operation or ensure
	// external durability. Do not block completion or introduce sensitive data.
	Observe(ctx context.Context, observation Observation)
}

// ObserverFunc adapts a function into an Observer.
type ObserverFunc func(ctx context.Context, observation Observation)

// Observe implements Observer.
func (f ObserverFunc) Observe(ctx context.Context, observation Observation) {
	f(ctx, observation)
}

func observeSafely(observer Observer, ctx context.Context, observation Observation) {
	defer func() { _ = recover() }()
	if observer != nil {
		observer.Observe(ctx, observation)
	}
}

func observationReason(err error) Reason {
	switch {
	case err == nil:
		return ReasonNone
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return ReasonCanceled
	case errors.Is(err, ErrBodyTooLarge), errors.Is(err, ErrSignatureHeadersTooLarge), errors.Is(err, ErrResponseTooLarge):
		return ReasonLimit
	case errors.Is(err, ErrInvalidTimestamp):
		return ReasonTimestamp
	case errors.Is(err, ErrReplay):
		return ReasonReplay
	case errors.Is(err, ErrReplayStore):
		return ReasonStore
	case errors.Is(err, ErrInvalidSignature), errors.Is(err, ErrMalformedSignatureHeader), errors.Is(err, ErrMalformedSignedHeader):
		return ReasonSignature
	case errors.Is(err, ErrEndpointRejected):
		return ReasonPolicy
	case errors.Is(err, ErrDeliveryFailed):
		return ReasonTransport
	default:
		return ReasonInternal
	}
}

func elapsed(clock func() time.Time, started time.Time) time.Duration {
	duration := clock().Sub(started)

	return max(duration, 0)
}
