# Changelog

All notable changes are documented here. This project follows Keep a Changelog
and Semantic Versioning.

## [Unreleased]

### Added

- Versioned HMAC-SHA-256 and HMAC-SHA-512 signing and verification.
- Signed, bounded nonces with injectable generation and a `crypto/rand` default.
- Exact-byte HTTP verification, bounded bodies and headers, rotation windows,
  safe typed failures, middleware, atomic replay protection, and a
  `go-idempotency` adapter.
- Deterministic envelopes, bounded delivery and retries, `Retry-After`,
  dead-letter and replay hooks, fan-out, SSRF and DNS-rebinding protection,
  and `go-queue` and `go-outbox` adapters.
- Secret-safe observations, independent Python vectors, fuzzers, allocation
  benchmarks, complete production coverage, and release gates.

### Fixed

- Saturate oversized numeric `Retry-After` values and reject body limits whose
  sentinel byte would overflow an `int64` bound.

## [1.0.0] - Unreleased

The first release will freeze the `v1` canonicalization and wire contracts.

[Unreleased]: https://github.com/faustbrian/go-webhook/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/faustbrian/go-webhook/releases/tag/v1.0.0
