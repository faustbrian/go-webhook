# Compatibility Policy

Each releasable directory is an independent Go module and follows semantic
versioning. Tags use `<module-directory>/v<version>`.

Before `v1`, minor releases MAY contain reviewed breaking changes, but every
break MUST be documented with migration guidance. Patch releases MUST remain
backward compatible. At and after `v1`, incompatible exported API or documented
behavior changes require a new major version.

Compatibility includes exported Go APIs, error classification, serialization,
protocol behavior, persistence schemas, environment variables, command output,
resource ownership, ordering, retry/idempotency semantics, and documented
defaults. A compile-compatible change can still be behaviorally breaking.

Version 2 changes the module path to `github.com/faustbrian/go-webhook/v2`
because the optional telemetry adapter accepts the v2 telemetry runtime.
See [migration guidance](docs/migration.md) for the import and runtime changes.

Version 3 was published as `github.com/faustbrian/go-webhook/v3` because the
public Outbox adapter uses separate Outbox v2 builder and envelope types.
This migration is not a compatible v2 patch. Version 3.0.1 preserves the
owned v3 API and delivery contracts; conditional dependency adoption
caveats are documented in the [integration guide](docs/integrations.md).

Specification-backed modules MUST NOT diverge from their declared standards.
Ambiguities require documented decisions and stable tests. The
[specification decision register](docs/specification-decisions.md) is part of
this compatibility contract. Deprecated APIs follow
[`DEPRECATION.md`](DEPRECATION.md).
