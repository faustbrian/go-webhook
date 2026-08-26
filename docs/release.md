# Release process

Releases start from a clean tree and the exact commit that will be tagged.
The changelog must describe every user-visible behavior and compatibility
decision.

## Required gates

Run the repository release contract documented in
[CONTRIBUTING.md](../CONTRIBUTING.md). It covers formatting, static analysis,
tests, exact production coverage, race detection, bounded fuzzing, benchmarks,
documentation, safety, interoperability, conformance, vulnerability scanning,
and module verification.

## Publication

Tag only the verified commit. The release workflow repeats the aggregate gate,
builds deterministic source archives, publishes checksums and provenance, and
verifies the resulting artifacts.

Provider claims must match the maintained provider matrix. A release is blocked
by failed required gates, a dirty source tree, unresolved high-severity
findings, or unsupported provider claims.

## Operational responsibility

Operators remain responsible for secret quality and rotation, durable-store
availability and tenant scoping, application idempotency and payload
validation, and any explicit SSRF-policy exceptions. The package does not claim
exactly-once delivery.
