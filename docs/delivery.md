# Delivery, retries, dead letters, and replay

`Deliverer` validates the endpoint before every attempt, generates stable
delivery and attempt IDs, signs the exact body and request target, bounds the
response, and records every completed attempt. Retryable transport failures
and HTTP 408, 425, 429, and 5xx responses use bounded exponential delay.
Delta-seconds and HTTP-date `Retry-After` values are honored up to the maximum.

Retries are forced to one attempt without an explicit idempotency key because
ambiguous receipt cannot otherwise be made safe. A queue or outbox consumer
uses `DeliverOnce`; the durable system remains the sole retry owner and nested
retry multiplication is prevented.

Terminal and exhausted deliveries invoke the optional dead-letter hook.
`Replay` invokes the audit hook with the original delivery ID and creates a
new delivery. It never reuses an attempt ID. Hook failures are returned and
hook panics cannot corrupt observation paths.

`FanOut` uses a fixed worker bound and preserves result order. It is an
orchestration primitive, not a durable queue. Cancellation stops new useful
work; callers own shutdown deadlines and persistence.
