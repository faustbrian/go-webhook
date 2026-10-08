package webhookotel

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	webhook "github.com/faustbrian/go-webhook/v3"
	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
)

func TestObserverRecordsBoundedMetricsAndCurrentSpanEvent(t *testing.T) {
	runtime, metrics, spans := exportingRuntime(t)
	observer, err := New(runtime)
	if err != nil {
		t.Fatal(err)
	}
	ctx, span := runtime.Tracer("webhook-contract").Start(context.Background(), "parent")
	event := webhook.Observation{Operation: webhook.OperationDeliveryAttempt, Outcome: webhook.OutcomeRetry, Reason: webhook.ReasonStatus, Duration: 250 * time.Millisecond, Algorithm: webhook.SHA256, StatusCode: 503, Classification: webhook.FailureRetryable}
	observer.Observe(ctx, event)
	span.End()
	event.Duration = 750 * time.Millisecond
	observer.Observe(context.Background(), event)
	flushContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := runtime.ForceFlush(flushContext); err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{"webhook.operation": "delivery_attempt", "webhook.outcome": "retry", "webhook.reason": "status", "webhook.algorithm": "sha256", "webhook.classification": "retryable", "http.response.status_class": "5xx"}
	snapshot := metrics.snapshot()
	if len(snapshot) != 2 {
		t.Fatalf("metric instruments = %+v", snapshot)
	}
	count := snapshot["webhook.operation.count"]
	if count.kind != "counter" || count.unit != "{operation}" || count.points != 1 || count.total != 2 || !reflect.DeepEqual(count.attributes, expected) {
		t.Fatalf("counter = %+v", count)
	}
	duration := snapshot["webhook.operation.duration"]
	if duration.kind != "histogram" || duration.unit != "s" || duration.points != 1 || duration.count != 2 || duration.total != 1 || !reflect.DeepEqual(duration.attributes, expected) {
		t.Fatalf("duration = %+v", duration)
	}
	ended := spans.GetSpans()
	if len(ended) != 1 || ended[0].Name != "parent" || len(ended[0].Events) != 1 {
		t.Fatalf("spans = %+v", ended)
	}
	eventRecord := ended[0].Events[0]
	if eventRecord.Name != "webhook.delivery_attempt" || !reflect.DeepEqual(attributeStrings(eventRecord.Attributes), expected) {
		t.Fatalf("span event = %+v", eventRecord)
	}
}

func TestNewAndStatusClassValidation(t *testing.T) {
	if _, err := New(nil); !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("New(nil) error = %v", err)
	}
	instrumentErr := errors.New("instrument unavailable")
	if _, err := newObserver(&errorMeter{Meter: metricnoop.NewMeterProvider().Meter("test"), err: instrumentErr}); !errors.Is(err, instrumentErr) {
		t.Fatalf("newObserver(counter) error = %v", err)
	}
	if _, err := newObserver(&errorMeter{Meter: metricnoop.NewMeterProvider().Meter("test"), err: instrumentErr, failHistogram: true}); !errors.Is(err, instrumentErr) {
		t.Fatalf("newObserver(histogram) error = %v", err)
	}
	for status, want := range map[int]string{0: "none", 99: "none", 100: "1xx", 200: "2xx", 503: "5xx", 599: "5xx", 600: "none"} {
		if got := statusClass(status); got != want {
			t.Fatalf("statusClass(%d) = %q, want %q", status, got, want)
		}
	}
}

type errorMeter struct {
	metric.Meter
	err           error
	failHistogram bool
}

func (m *errorMeter) Int64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	if !m.failHistogram {
		return nil, m.err
	}
	return m.Meter.Int64Counter(name, options...)
}

func (m *errorMeter) Float64Histogram(string, ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
	return nil, m.err
}
