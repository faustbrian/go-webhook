package webhookoutbox_test

import (
	"context"

	outbox "github.com/faustbrian/go-transactional-outbox/v2"
	webhook "github.com/faustbrian/go-webhook/v3"
	webhookoutbox "github.com/faustbrian/go-webhook/v3/adapters/outbox"
)

// The adapter's public signatures must compose with the released producer
// types, not an earlier nominally distinct envelope or builder.
var (
	_ interface {
		Publish(context.Context, outbox.Envelope) error
	} = (*webhookoutbox.Publisher)(nil)
	_ func(*outbox.EnvelopeBuilder, string, webhook.DeliveryRequest, int) (outbox.Envelope, error) = webhookoutbox.Build
)
