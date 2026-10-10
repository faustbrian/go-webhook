# Integrations and observability

## HTTP clients

`Deliverer` accepts any `HTTPDoer`. The released
[`http-client`](https://pkg.go.dev/github.com/faustbrian/go-http-client)
`Client` implements that interface without making it a required dependency.
When substituting it for `NewSecureHTTPClient`, configure its egress, proxy,
DNS, redirect, timeout, and response-ownership policies so they preserve the
webhook delivery guarantees.

## Queue and outbox

`adapters/queue` encodes a bounded versioned delivery and defaults queue
retry count to zero. Its handler performs one attempt. `adapters/outbox`
builds an outbox envelope and its publisher likewise performs one attempt.
Durable settlement, visibility, retries, and dead lettering remain owned by
those systems.

The selected `go-queue` v1.1.3 preserves this adapter’s consumed `core` and
`job` contracts. Applications that also imported
`github.com/faustbrian/go-queue/rabbitmq` from the older root module must add
the separately published module at v1.0.2. This is not a drop-in runtime
migration: `WithNativeConfig`, stable message identity, and explicit native
connection and producer/consumer policies are required. Automatic
acknowledgement, fanout, and headers exchanges are rejected; topology remains
application-owned. The producer opens during construction and the consumer
opens lazily. See the published [RabbitMQ migration notes](https://github.com/faustbrian/go-queue/tree/376b2ea5c2bd25ace6af4374e464852203df1f8c/rabbitmq/CHANGELOG.md#L38-L59)
and [configuration guidance](https://github.com/faustbrian/go-queue/tree/376b2ea5c2bd25ace6af4374e464852203df1f8c/rabbitmq/README.md#L79-L86).

## Telemetry and logs

`Observer` emits only fixed operation, outcome, reason, algorithm,
classification, attempt, bounded status, and duration fields. Adapt it to a
metric or trace backend without adding event ID, URL, host, query, payload,
signature, key ID, header, or raw error attributes.

For `telemetry/v2`, pass its initialized runtime to
`adapters/otel.New`. The adapter records fixed metrics and adds an event
to the current span. `InstrumentHTTPClient` clones an existing explicit client
and wraps its transport with tracing, metrics, and W3C propagation; use the
client returned by `NewSecureHTTPClient` so SSRF and redirect policy remain the
base transport. `adapters/slog.New` accepts the `*slog.Logger` returned by
`log` and emits only the fixed observation schema. Its redaction handler is
still recommended as defense in depth. Disabled telemetry is a nil `Observer`
and has no business-path dependency.

Version 2 telemetry disables traces, metrics, and global registration by
default. Explicitly enable the signals the application needs. The webhook
adapter uses the supplied runtime directly, so global registration is not
required for it; enable registration only if other application code relies on
OpenTelemetry globals. Plaintext Collector transport now requires an explicit
opt-in and should be used only within a protected trust zone.

The selected OpenTelemetry SDK v1.47.0 [removes the experimental export
batch-size setting](https://github.com/open-telemetry/opentelemetry-go/tree/66cfc9520e205b7d450183532772401bc2b6674c/CHANGELOG.md#L70)
`OTEL_GO_X_METRIC_EXPORT_BATCH_SIZE`. Applications relying on that variable
lose their configured export batch-size limit. The released `telemetry/v2`
v2.0.0 runtime does not expose the reader option that replaces it, and this
webhook adapter cannot configure that limit. Exporter modules remain at
v1.45.0; this SDK update does not include later exporter fixes.

## Deterministic consumer tests

The `webhooktest` package creates a signer/verifier pair with a manual clock,
sequential signed nonces, and sequential delivery IDs. It has no global state
and is concurrency-safe. Use it for repeatable receiver, sender, retry, and
rotation tests without weakening production entropy.
