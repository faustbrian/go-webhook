package webhookotel

import (
	"context"
	"sync"
	"testing"
	"time"

	telemetry "github.com/faustbrian/go-telemetry/v2"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type metricSnapshot struct {
	kind, unit string
	points     int
	count      uint64
	total      float64
	attributes map[string]string
}
type contractMetricExporter struct {
	mu      sync.Mutex
	metrics map[string]metricSnapshot
}

func (*contractMetricExporter) Temporality(kind sdkmetric.InstrumentKind) metricdata.Temporality {
	return sdkmetric.DefaultTemporalitySelector(kind)
}
func (*contractMetricExporter) Aggregation(kind sdkmetric.InstrumentKind) sdkmetric.Aggregation {
	return sdkmetric.DefaultAggregationSelector(kind)
}
func (e *contractMetricExporter) Export(_ context.Context, resource *metricdata.ResourceMetrics) error {
	// Exported SDK buffers may be reused immediately. Retain owned scalar and
	// attribute copies, never ResourceMetrics or its slices.
	captured := make(map[string]metricSnapshot)
	for _, scope := range resource.ScopeMetrics {
		if scope.Scope.Name != "github.com/faustbrian/go-webhook/v3/adapters/otel" {
			continue
		}
		for _, metric := range scope.Metrics {
			item := metricSnapshot{unit: metric.Unit}
			switch data := metric.Data.(type) {
			case metricdata.Sum[int64]:
				item.kind = "counter"
				item.points = len(data.DataPoints)
				if item.points == 1 {
					item.total = float64(data.DataPoints[0].Value)
					item.attributes = attributeStrings(data.DataPoints[0].Attributes.ToSlice())
				}
			case metricdata.Histogram[float64]:
				item.kind = "histogram"
				item.points = len(data.DataPoints)
				if item.points == 1 {
					item.total = data.DataPoints[0].Sum
					item.count = data.DataPoints[0].Count
					item.attributes = attributeStrings(data.DataPoints[0].Attributes.ToSlice())
				}
			default:
				item.kind = "unexpected"
			}
			captured[metric.Name] = item
		}
	}
	e.mu.Lock()
	e.metrics = captured
	e.mu.Unlock()
	return nil
}
func (*contractMetricExporter) ForceFlush(context.Context) error { return nil }
func (*contractMetricExporter) Shutdown(context.Context) error   { return nil }
func (e *contractMetricExporter) snapshot() map[string]metricSnapshot {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.metrics
}
func attributeStrings(attributes []attribute.KeyValue) map[string]string {
	values := make(map[string]string, len(attributes))
	for _, attribute := range attributes {
		values[string(attribute.Key)] = attribute.Value.AsString()
	}
	return values
}
func exportingRuntime(t *testing.T) (*telemetry.Runtime, *contractMetricExporter, *tracetest.InMemoryExporter) {
	t.Helper()
	metrics := &contractMetricExporter{}
	spans := tracetest.NewInMemoryExporter()
	config := telemetry.DefaultConfig("webhook-contract", "v1")
	config.Traces.Enabled = true
	config.Traces.Sampler.Ratio = 1
	config.Metrics.Enabled = true
	runtime, err := telemetry.Init(context.Background(), config, telemetry.WithTraceExporter(spans), telemetry.WithMetricExporter(metrics))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := runtime.Shutdown(ctx); err != nil {
			t.Errorf("Shutdown: %v", err)
		}
	})
	return runtime, metrics, spans
}
