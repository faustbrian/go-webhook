package webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ErrInvalidEnvelope identifies missing required fields or invalid JSON data.
var ErrInvalidEnvelope = errors.New("invalid webhook envelope")

// Envelope is a small protocol-independent event envelope. Its v1 JSON field
// order and UTC timestamp representation are compatibility-sensitive. This is
// local WEBHOOK-DEC-013, not CloudEvents 1.0.2 conformance. Data/maps are borrowed;
// callers own size bounds and concurrent mutation. Zero Envelope is invalid.
type Envelope struct {
	// ID must be nonempty local event identity.
	ID string
	// Type must be nonempty application event type.
	Type string
	// Source must be nonempty application source; no URI validation is added.
	Source string
	// Subject is optional application data omitted when empty.
	Subject string
	// Time must be nonzero and emits UTC RFC3339Nano.
	Time time.Time
	// Data is borrowed valid JSON, including null; not redacted or size-bounded here.
	Data json.RawMessage
	// Metadata is borrowed arbitrary string data emitted as a nested object; nil/empty omits it.
	Metadata map[string]string
}

// MarshalJSON validates required fields/data and emits owned deterministic JSON
// using Go encoding/json, UTC RFC3339Nano and application/json. Embedded JSON is
// compacted/escaped by the encoder, not byte-for-byte preserved. No application
// schema, secret redaction or output allocation bound is added.
func (e Envelope) MarshalJSON() ([]byte, error) {
	if e.ID == "" || e.Type == "" || e.Source == "" || e.Time.IsZero() || !json.Valid(e.Data) {
		return nil, fmt.Errorf("%w: id, type, source, time, and valid JSON data are required", ErrInvalidEnvelope)
	}
	wire := struct {
		SpecVersion     string            `json:"specversion"`
		ID              string            `json:"id"`
		Type            string            `json:"type"`
		Source          string            `json:"source"`
		Subject         string            `json:"subject,omitempty"`
		Time            string            `json:"time"`
		DataContentType string            `json:"datacontenttype"`
		Data            json.RawMessage   `json:"data"`
		Metadata        map[string]string `json:"metadata,omitempty"`
	}{
		SpecVersion:     "1.0",
		ID:              e.ID,
		Type:            e.Type,
		Source:          e.Source,
		Subject:         e.Subject,
		Time:            e.Time.UTC().Format(time.RFC3339Nano),
		DataContentType: "application/json",
		Data:            e.Data,
		Metadata:        e.Metadata,
	}

	return json.Marshal(wire)
}
