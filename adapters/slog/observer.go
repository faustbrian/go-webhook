// Package webhookslog adapts secret-safe webhook observations to log's standard
// slog surface. New borrows a logger and Observe uses a fixed attribute schema.
// No handler, transport, redaction or logger lifecycle is owned. Privacy depends
// on core-generated categories or a trusted direct caller: string-backed enums
// are rendered unchanged and prebound/encoder data remain outside this adapter.
package webhookslog

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	webhook "github.com/faustbrian/go-webhook/v3"
)

// ErrInvalidConfig means no logger was supplied.
var ErrInvalidConfig = errors.New("webhook/slog: logger is required")

// Observer emits a fixed structured record. Its schema intentionally has no
// payload, URL, header, event ID, key ID, signature, nonce, or raw error field.
type Observer struct {
	logger *slog.Logger
	level  slog.Level
}

var _ webhook.Observer = (*Observer)(nil)

// New constructs a log-compatible observer. log returns *slog.Logger,
// so no logger-specific interface or global state is required.
func New(logger *slog.Logger, level slog.Level) (*Observer, error) {
	if logger == nil {
		return nil, ErrInvalidConfig
	}

	return &Observer{logger: logger, level: level}, nil
}

// Observe emits one synchronous record at configured level. Direct event values
// are neither bounded nor normalized; logger handler errors follow slog behavior.
// Call with non-nil ctx/constructed Observer and use a concurrent-safe logger
// pipeline. Core invocation contains observer panics, direct invocation does not.
func (o *Observer) Observe(ctx context.Context, event webhook.Observation) {
	o.logger.LogAttrs(ctx, o.level, "webhook lifecycle",
		slog.String("webhook.operation", string(event.Operation)),
		slog.String("webhook.outcome", string(event.Outcome)),
		slog.String("webhook.reason", string(event.Reason)),
		slog.String("webhook.algorithm", string(event.Algorithm)),
		slog.String("webhook.classification", string(event.Classification)),
		slog.String("http.response.status_class", statusClass(event.StatusCode)),
		slog.Int("webhook.attempt", event.Attempt),
		slog.Int64("webhook.duration_ms", event.Duration.Milliseconds()),
	)
}

func statusClass(status int) string {
	if status < 100 || status > 599 {
		return "none"
	}

	return strconv.Itoa(status/100) + "xx"
}
