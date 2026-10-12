# Source documentation audit

## Scope and contract inventory

Baseline `172a21507211a53d5f775eca5fd471dfb824e1b6` contains 45 tracked Go
files (16 maintained source, 29 test) in one v3 module and seven packages:
root webhook, five adapters and webhooktest. Reproduce with
`git ls-files '*.go' 'go.mod' '**/go.mod'` and
`golib source-docs inventory --json`. All tracked code/tests, six fuzz targets,
two benchmarks and five executable examples were inspected. There are no
tracked generated/vendor, build-tag/platform or embedded-template variants,
nor generation, deprecation or actionable debt directives. The baseline has
171 objective declaration/member-comment gaps.

| Owner | Audited contract and documentation need |
| --- | --- |
| `signature.go` | Local v1 canonical form, one algorithm/nonce, copied key material, inclusive signed-time rotation, candidate ordering, constant-time MAC comparison versus variable total runtime, replay admission versus business effects. |
| `http.go`, `body.go`, `middleware.go` | Early header limits, exact remaining stream, body closure/restoration and partial failure returns, borrowed metadata/extractors, authenticated-data ID provenance and internal error exposure. |
| `delivery.go`, `delivery_wire.go` | Explicit retry/idempotency policy, zero-or-one queue/outbox attempt, response ownership/sensitivity, per-call bounds, callbacks/cancellation, joined dead-letter errors, fan-out join and sensitive durable wire data. |
| `ssrf.go` | Copied allow/deny prefixes, resolver obligations, all-answer admission, validated dial targets, no proxy/redirect/default HTTP2 and caller-owned client lifecycle. |
| `envelope.go`, `observation.go` | Local JSON profile versus CloudEvents, raw application data/limits, fixed core observations versus arbitrary direct caller values. |
| Adapters | Idempotency lease availability/expiry, queue/outbox retry and settlement ownership, borrowed options, telemetry/logging lifecycle and direct-observation trust. |
| `webhooktest` and test helpers | Deterministic but non-production nonce/ID sources, lock-order sequencing, synthetic lease/egress/exporter fixtures and cleanup ownership. |

Pinned authority and interpretations remain in
[specification decisions](specification-decisions.md) and
[the source manifest](../specification/manifest.tsv). New adjacent comments
reference specific decisions/tests instead of inventing interoperability claims.
No provider, RFC 9421 or CloudEvents conformance is inferred from local v1.

## Corrected current decisions

`WEBHOOK-DEC-012` previously promised UTF-8 validation in `HeaderEventID`, but
the maintained function only checks a single nonempty value and its byte length.
It also does not automatically sign an arbitrary selected header. Current
selected-behavior prose/JSON and its rendered mirrors now describe that actual
boundary and application-owned binding. Preserve historical native digests and
append the corrected typed digests; no runtime tightening or external source-pin
change is included. Tests already protect extraction order/count/size, not
UTF-8 validation or arbitrary-header MAC coverage. `WEBHOOK-DEC-016` and the
queue/outbox reference prose now say at most one attempt: admission failures
can return without any HTTP send.

## Selected proof and preservation

Preserve executable tokens/literals, directives, example outputs, protocol
vectors, historical specification evidence, dependency/configuration bytes and
API signatures. Only DEC-012/016 current prose and their append-only digests change.
Native API mode is `apidiff` at `api/baseline-v3.txt`, not full Go-doc prose;
current and historical symbol baselines remain unchanged. The native docs script
renders all seven packages and runs existing examples; the full public render
and canonical package comments receive ordinary inspection. All tracked Go
files receive parser/token/directive checks independently of package discovery.
Local links, changed-content secrets and the complete final diff are inspected.

No unchanged cryptographic conformance, real DNS/egress, service, active fuzz,
race, performance or release campaigns are selected for Tier A prose. Existing
example TestMain cleanup remains enabled. Source proof does not qualify runtime
security, publication or CI; the coordinator owns main delivery and shared
objective enforcement. No release or repin is required by these comments.

## Objective Enforcement

The qualified public Tools v2.4.0 binary and immutable reusable CI source
provide the strict tracked-source check through `make source-docs-check`,
the existing local check/CI routes and Required aggregation. Release and
selected security/assurance routes retain their existing controls. This
structural status supplements the rendered semantic audit; it does not
certify prose meaning, generator execution or runtime conformance.
