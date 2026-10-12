// Package webhookoutbox adapts webhook deliveries to outbox envelopes and relay
// publishers. Build only creates an envelope; callers insert it transactionally.
// NewPublisher/Publish borrow a Deliverer and cap core attempts at one, while
// the relay owns claiming, retries, dead-letter policy and settlement. No database,
// background relay, resource shutdown or credential policy is installed here.
package webhookoutbox

import (
	"context"
	"errors"
	"fmt"

	"github.com/faustbrian/go-transactional-outbox/v2"
	webhook "github.com/faustbrian/go-webhook/v3"
)

// ErrInvalidConfig identifies missing required builder/deliverer/topic or payload bound.
var ErrInvalidConfig = errors.New("webhook/outbox: invalid configuration")

// Build encodes a bounded delivery request into a validated outbox envelope.
// Payload contains unredacted application data. Builder callbacks/lifetime stay
// caller-owned; no store insertion or transaction is performed. EventID becomes
// OrderingKey, IdempotencyKey remains explicit, PayloadVersion is one.
func Build(
	builder *outbox.EnvelopeBuilder,
	topic string,
	delivery webhook.DeliveryRequest,
	maxMessageBytes int,
) (outbox.Envelope, error) {
	if builder == nil || topic == "" || maxMessageBytes <= 0 {
		return outbox.Envelope{}, ErrInvalidConfig
	}
	payload, err := webhook.MarshalDeliveryRequest(delivery, maxMessageBytes)
	if err != nil {
		return outbox.Envelope{}, err
	}
	envelope, err := builder.Build(outbox.NewEnvelopeParams{
		Topic:          topic,
		Payload:        payload,
		PayloadVersion: 1,
		Metadata:       delivery.Metadata,
		OrderingKey:    delivery.EventID,
		IdempotencyKey: delivery.IdempotencyKey,
	})
	if err != nil {
		return outbox.Envelope{}, fmt.Errorf("webhook/outbox: build envelope: %w", err)
	}

	return envelope, nil
}

// Publisher implements relay.Publisher with at most one HTTP attempt per
// claimed envelope. The relay owns retries/settlement; a supplied Deliverer's
// DeadLetter callback can still run on its terminal path, so coordinate ownership
// explicitly. Zero Publisher is unusable; shared calls require a safe Deliverer.
type Publisher struct {
	deliverer *webhook.Deliverer
	maxBytes  int
}

// NewPublisher validates the deliverer and payload bound.
func NewPublisher(deliverer *webhook.Deliverer, maxMessageBytes int) (*Publisher, error) {
	if deliverer == nil || maxMessageBytes <= 0 {
		return nil, ErrInvalidConfig
	}

	return &Publisher{deliverer: deliverer, maxBytes: maxMessageBytes}, nil
}

// Publish decodes bounded payload and requests at most one HTTP send. It does
// not validate envelope metadata/topic/version or settle the claim; relay callers
// own those checks and durable outcomes. Delivery errors remain wrapped.
func (p *Publisher) Publish(ctx context.Context, envelope outbox.Envelope) error {
	delivery, err := webhook.UnmarshalDeliveryRequest(envelope.Payload, p.maxBytes)
	if err != nil {
		return err
	}
	if _, err := p.deliverer.DeliverOnce(ctx, delivery); err != nil {
		return fmt.Errorf("webhook/outbox: publish delivery: %w", err)
	}

	return nil
}
