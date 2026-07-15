# Integrations and observability

## HTTP clients

`Deliverer` accepts any `HTTPDoer`. The named `go-http-client` repository has
no published Go module or API as of 2026-07-15, so no concrete import can be
truthfully compiled. Its eventual client can integrate by implementing
`Do(*http.Request) (*http.Response, error)`. It must preserve the SSRF policy,
no-proxy transport, dial-time validation, redirect refusal, and response
ownership supplied by `NewSecureHTTPClient`.

## Queue and outbox

`adapters/goqueue` encodes a bounded versioned delivery and defaults queue
retry count to zero. Its handler performs one attempt. `adapters/gooutbox`
builds an outbox envelope and its publisher likewise performs one attempt.
Durable settlement, visibility, retries, and dead lettering remain owned by
those systems.

## Telemetry and logs

`Observer` emits only fixed operation, outcome, reason, algorithm,
classification, attempt, bounded status, and duration fields. Adapt it to a
metric or trace backend without adding event ID, URL, host, query, payload,
signature, key ID, header, or raw error attributes.

For `go-telemetry`, pass its initialized runtime to
`adapters/gotelemetry.New`. The adapter records fixed metrics and adds an event
to the current span. Its HTTP instrumentation may wrap the secure policy
transport but must not replace it. For `go-log`, map observations to fixed
`slog.Attr` values and use its redaction handler as defense in depth. Disabled
telemetry is a nil `Observer` and has no business-path dependency.
