package webhookotel

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	telemetry "github.com/faustbrian/go-telemetry/v2"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TestInstrumentHTTPClientInjectsTraceAndPreservesClientPolicy(t *testing.T) {
	runtime, _, spans := exportingRuntime(t)
	base := &recordingTransport{}
	redirects := 0
	redirect := func(*http.Request, []*http.Request) error { redirects++; return http.ErrUseLastResponse }
	original := &http.Client{Transport: base, Timeout: time.Second, CheckRedirect: redirect}
	client, err := InstrumentHTTPClient(runtime, original, "webhook.deliver")
	if err != nil {
		t.Fatal(err)
	}
	if client == original || original.Transport != base || client.Timeout != original.Timeout || client.CheckRedirect == nil {
		t.Fatal("caller client policy changed")
	}
	parent := trace.NewSpanContext(trace.SpanContextConfig{TraceID: trace.TraceID{1}, SpanID: trace.SpanID{2}, TraceFlags: trace.FlagsSampled})
	ctx, cancel := context.WithTimeout(trace.ContextWithSpanContext(context.Background(), parent), 5*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://example.com/hook", nil)
	if err != nil {
		t.Fatal(err)
	}
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	if err := response.Body.Close(); err != nil {
		t.Fatal(err)
	}
	extracted := propagation.TraceContext{}.Extract(context.Background(), propagation.MapCarrier{"traceparent": base.traceparent})
	received := trace.SpanContextFromContext(extracted)
	if !received.IsValid() || received.TraceID() != parent.TraceID() || !received.IsSampled() || base.calls != 1 {
		t.Fatalf("trace propagation = %s calls=%d", base.traceparent, base.calls)
	}
	base.status = http.StatusFound
	response, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	if err := response.Body.Close(); err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusFound || redirects != 1 || base.calls != 2 {
		t.Fatalf("redirect policy status=%d redirects=%d sends=%d", response.StatusCode, redirects, base.calls)
	}
	denied := errors.New("caller transport refused egress")
	base.err = denied
	if _, err := client.Do(request); !errors.Is(err, denied) {
		t.Fatalf("transport rejection = %v", err)
	}
	if base.calls != 3 || original.Transport != base {
		t.Fatalf("caller transport sends=%d", base.calls)
	}
	if err := runtime.ForceFlush(ctx); err != nil {
		t.Fatal(err)
	}
	ended := spans.GetSpans()
	if len(ended) != 3 {
		t.Fatalf("exported HTTP spans = %+v", ended)
	}
	for _, span := range ended {
		if span.SpanContext.TraceID() != parent.TraceID() || span.Parent.SpanID() != parent.SpanID() {
			t.Fatalf("HTTP span lost parent = %+v", span)
		}
	}
}

func TestInstrumentHTTPClientRejectsUnsafeConfiguration(t *testing.T) {
	runtime := testRuntime(t)
	base := &http.Client{Transport: &recordingTransport{}}
	for name, test := range map[string]struct {
		runtime *telemetry.Runtime
		client  *http.Client
		op      string
	}{
		"runtime":   {client: base, op: "webhook.deliver"},
		"client":    {runtime: runtime, op: "webhook.deliver"},
		"transport": {runtime: runtime, client: &http.Client{}, op: "webhook.deliver"},
		"operation": {runtime: runtime, client: base, op: "bad operation"},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := InstrumentHTTPClient(test.runtime, test.client, test.op); err == nil {
				t.Fatal("InstrumentHTTPClient() accepted unsafe configuration")
			}
		})
	}
}

func testRuntime(t *testing.T) *telemetry.Runtime {
	t.Helper()
	config := telemetry.DefaultConfig("webhook-test", "v1")
	config.Traces.Enabled = false
	config.Metrics.Enabled = false
	runtime, err := telemetry.Init(context.Background(), config)
	if err != nil {
		t.Fatalf("telemetry.Init() error = %v", err)
	}
	t.Cleanup(func() {
		if err := runtime.Shutdown(context.Background()); err != nil {
			t.Errorf("Shutdown() error = %v", err)
		}
	})
	return runtime
}

type recordingTransport struct {
	traceparent string
	calls       int
	status      int
	err         error
}

func (t *recordingTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if request == nil {
		return nil, errors.New("request required")
	}
	t.calls++
	t.traceparent = request.Header.Get("Traceparent")
	if t.err != nil {
		return nil, t.err
	}
	status := t.status
	if status == 0 {
		status = http.StatusNoContent
	}
	header := make(http.Header)
	if status == http.StatusFound {
		header.Set("Location", "https://example.com/not-followed")
	}
	return &http.Response{StatusCode: status, Header: header, Body: io.NopCloser(strings.NewReader(""))}, nil
}
