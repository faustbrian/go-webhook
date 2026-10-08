# Migration and SemVer

<a id="version-3-module-migration-pending-publication"></a>

## Version 3 module migration

The published v3 line uses `github.com/faustbrian/go-webhook/v3`. Insert `/v3` before each
webhook package suffix and use `github.com/faustbrian/go-transactional-outbox/v2`
for builders, envelopes and relay contracts passed to `adapters/outbox`.
The adapter keeps delivery encoding, ownership, cancellation and one HTTP
attempt per publish; the relay remains responsible for retries and dead letters.
The telemetry instrumentation scope changes to the `/v3` package identity.
No signing, verification or delivery-policy change is intended by this migration.
Version 3.0.0 was published on 2026-10-08. For later v3 dependency
updates, also review the [integration caveats](integrations.md).

The Idempotency ecosystem compatibility harness and Service external-reference
harness still consume historical webhook v1. They are future migration targets,
not current v3 reverse dependencies or evidence of v3 adoption.

## Version 2 module migration

Change all `github.com/faustbrian/go-webhook` package imports to
`github.com/faustbrian/go-webhook/v2`. The optional `adapters/otel` entry points
now take `*telemetry.Runtime` from
`github.com/faustbrian/go-telemetry/v2`; migrate the telemetry runtime import
and construction at the same time. The webhook signing, verification, replay,
delivery, and observation contracts are otherwise unchanged.

Telemetry v2 does not implicitly enable traces, metrics, or global
registration. Configure each signal your application requires explicitly;
the webhook adapter does not require global registration because it receives
the runtime directly. For protected plaintext Collector connections, set the
trace and metric exporter TLS `Insecure` options explicitly. See the
[telemetry v2 migration contract](https://github.com/faustbrian/go-telemetry/blob/v2.0.0/docs/compatibility.md#version-2-migration).

There is no legacy API before v1. Adoption should first deploy verification in
shadow observation mode using synthetic fixtures, then enforce signatures,
then enable a tenant-scoped replay store, and finally enable outbound retries
only after endpoint idempotency is proven.

Changing canonical fields, nonce handling, encodings, ordering, line endings,
header grammar, body digest, envelope wire bytes, exported error identities,
retryable status classification, or an existing provider preset requires a
major version.
Adding an isolated algorithm or provider can be minor when negotiation cannot
downgrade existing behavior. Security fixes may intentionally reject input
that was previously accepted and will be called out in the changelog.

Every protocol change must update the affected stable entry in the
[specification decision register](specification-decisions.md). Superseded
entries remain linked so adopters can distinguish an intentional migration
from an undocumented interpretation change.
