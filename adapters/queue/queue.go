// Package webhookqueue adapts bounded delivery requests to queue messages.
// New/Enqueue borrow a producer; Handle decodes and caps core attempts at one.
// The external queue owns durability, retries, settlement and shutdown. No
// consumer worker, transport or redaction is installed by this adapter.
package webhookqueue

import (
	"context"
	"errors"
	"fmt"

	"github.com/faustbrian/go-queue/core"
	"github.com/faustbrian/go-queue/job"
	webhook "github.com/faustbrian/go-webhook/v3"
)

// ErrInvalidConfig identifies missing queue/deliverer or nonpositive payload bound.
var ErrInvalidConfig = errors.New("webhook/queue: invalid configuration")

// Queue is the narrow synchronous producer surface implemented by queue.
type Queue interface {
	// Queue synchronously accepts/rejects a message and options. It has no ctx:
	// bound its blocking and support concurrent use if adapters are shared.
	// Message.Bytes returns a fresh copy; options/pointer fields remain borrowed.
	Queue(message core.QueuedMessage, options ...job.AllowOption) error
}

// Config bounds messages and optionally assigns queue-level retry policy.
// Delivery handling performs at most one HTTP attempt, preventing nested retry
// multiplication.
type Config struct {
	// Queue is mandatory and borrowed; typed nil is not detected.
	Queue Queue
	// MaxMessageBytes must be positive and caps encoded delivery payload, not queue storage.
	MaxMessageBytes int
	// JobOptions is slice-copied; pointer fields remain shared. Empty defaults to RetryCount zero.
	JobOptions []job.AllowOption
}

// Adapter enqueues versioned delivery requests. Construct with New; zero is
// unusable. Shared use requires a safe queue and immutable borrowed option
// pointer values. It owns no queue resource lifecycle or global state.
type Adapter struct {
	queue      Queue
	maxBytes   int
	jobOptions []job.AllowOption
}

// New validates the queue and message bound.
func New(config Config) (*Adapter, error) {
	if config.Queue == nil || config.MaxMessageBytes <= 0 {
		return nil, ErrInvalidConfig
	}
	options := append([]job.AllowOption(nil), config.JobOptions...)
	if len(options) == 0 {
		options = []job.AllowOption{{RetryCount: job.Int64(0)}}
	}

	return &Adapter{queue: config.Queue, maxBytes: config.MaxMessageBytes, jobOptions: options}, nil
}

// Enqueue checks cancellation, encodes a bounded request, and synchronously
// hands it to queue. Acceptance is not an exactly-once guarantee. ctx must be
// non-nil; cancellation is checked only before encoding, not propagated into the
// contextless Queue call. Raw serialized data may be sensitive; queue errors are
// wrapped verbatim and callback panics propagate.
func (a *Adapter) Enqueue(ctx context.Context, delivery webhook.DeliveryRequest) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	encoded, err := webhook.MarshalDeliveryRequest(delivery, a.maxBytes)
	if err != nil {
		return err
	}
	if err := a.queue.Queue(message(encoded), a.jobOptions...); err != nil {
		return fmt.Errorf("webhook/queue: enqueue delivery: %w", err)
	}

	return nil
}

// Handle strictly decodes one message and requests at most one HTTP attempt. The
// surrounding queue owns any retry schedule and settlement behavior. Admission
// failures may perform zero HTTP sends. No queue acknowledgement/release occurs;
// application persistence and retry classification remain caller-owned.
func Handle(
	ctx context.Context,
	deliverer *webhook.Deliverer,
	encoded []byte,
	maxMessageBytes int,
) (webhook.DeliveryResult, error) {
	if deliverer == nil || maxMessageBytes <= 0 {
		return webhook.DeliveryResult{}, ErrInvalidConfig
	}
	delivery, err := webhook.UnmarshalDeliveryRequest(encoded, maxMessageBytes)
	if err != nil {
		return webhook.DeliveryResult{}, err
	}

	return deliverer.DeliverOnce(ctx, delivery)
}

type message []byte

func (m message) Bytes() []byte {
	return append([]byte(nil), m...)
}
