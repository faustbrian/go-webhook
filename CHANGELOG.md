# Changelog

All notable changes are documented here. This project follows Keep a Changelog
and Semantic Versioning.

## [Unreleased]

### Documentation

- Remove the archived monorepo documentation link; package guidance remains in
  the repository-owned documentation.

### Added

- Make the [specification decision register](docs/specification-decisions.md)
  machine-auditable with pinned source and change authorities, explicit
  conformance and differential evidence, canonical history, and fail-closed
  change control. Current decision content:
  - WEBHOOK-DEC-001 sha256:20fb87d9166b3063bb113c8273b9bd688f67e7f3deba6d28fa127445572c3c19
  - WEBHOOK-DEC-002 sha256:430ae1045ec2d8ea27c6de38d60077d5290c813b66107af4db9e9dce72cf627e
  - WEBHOOK-DEC-003 sha256:5561a06a38fa997561207df636d4ab3398ce4c7b0b5e108f295c1ef88d24e16b
  - WEBHOOK-DEC-004 sha256:8a4d6659284fdbc89aa977b2c294ff50f41b9687da5c169c189e9dc98d08b77e
  - WEBHOOK-DEC-005 sha256:8c6518536c7817b373c90ae63a655588f2d06683af4a7da05e3a4fa30f08e9f7
  - WEBHOOK-DEC-006 sha256:8f9075521ae7fa19362d63c88a45ef454802716ade0446572fbe16421c7770d2
  - WEBHOOK-DEC-007 sha256:e03ba74b8615f3cace63f39f3f9fcb1321c779691ffb2f9340f636c8bb7ad5c8
  - WEBHOOK-DEC-008 sha256:70a6899937062bf17d5161d2a6821aaab0e9fb7ac14b225e08340630702ea3a3
  - WEBHOOK-DEC-009 sha256:296686e7b21b2ab30e676748baabce7dff4461b4d9b34551f7d6a6cdea413524
  - WEBHOOK-DEC-010 sha256:9412aee18a5c64c3c53543f222bdbb1abe80b180e0cec6eb5237627c967c95cd
  - WEBHOOK-DEC-011 sha256:cddd7b2bfe8a5a81275c15311519ef4f2842bb4b320eadff45173a04d2271f01
  - WEBHOOK-DEC-012 sha256:aa3beea1fd172bfa11c16f930f5e32cfde07fd0e4e920d3cd4c5e6f5b7d1af4f
  - WEBHOOK-DEC-013 sha256:dee4801524b0f81b02d89ff50ab9cac4c4186ab3912b0f700c77c1b5918a1d09
  - WEBHOOK-DEC-014 sha256:03ce7aa6765ce943feadd178692deff0640e56945d9000012ea875064a5a71ee
  - WEBHOOK-DEC-015 sha256:aef99c818eb2b6ab32410d98611da8c6feb56f7941ae50c73bb9308e1dbf6a65
  - WEBHOOK-DEC-016 sha256:634c71628df1d835acbd35b7259cbf49f441166494ac2449730f2c35d5c80b42
  - WEBHOOK-DEC-017 sha256:d0c06d8f8a69672ac081dce3cf785bdec6e7a02e47c536cf6a1108a289cf3ab7
  - WEBHOOK-DEC-018 sha256:0641b615cf0ecfa2375d10c2e173bb4296394774dfdf027698bb040a1889fdc5
  - WEBHOOK-DEC-019 sha256:76234172ed1b5d3c4721d7942f2e3c31e4c0651cb7ce9897011f1194b060ee44
  - WEBHOOK-DEC-020 sha256:bd28955089597f1311fc1c926fb1a281eddd740a0fecc444583a4e395a625c72
  - WEBHOOK-DEC-021 sha256:ec64861d6fe76a27da80add722eccc487c0612cc0565711e7d4ca70f8544580b
  - WEBHOOK-DEC-022 sha256:189d07fd4c2962dda1b5994e6d5a2622697e660de45df11beb0e79a87bab1540
  - WEBHOOK-DEC-023 sha256:bcad644089ba17d7747a815e61ab723a055ef54811c820a4691698c20713d34f
  - WEBHOOK-DEC-024 sha256:d634389967f7a938b54ba5434d308ce71ab0334fc35990051e8b5accbf50b717
  - WEBHOOK-DEC-025 sha256:36033f63cabaf3dadd61c00bb7d521cdbed564fbc20aec17ac50d24f2f79cd53
  - WEBHOOK-DEC-026 sha256:6fb346a86a7f4b5c0244f921c008c127e512189da03f60e9f8b2a047d4ff6a1a

## [1.0.0] - 2026-08-25

### Compatibility

- Regenerate the exported API baseline with the repository's Go 1.26
  toolchain so the stable `Envelope.Data` contract is represented accurately.

### Changed

- Exclude intentional nested modules from root local-proxy archives so local,
  bootstrap, CI, and public module checksums describe the same source
  boundary.

- Track the pinned documentation-tool lockfile so clean CI checkouts install
  the exact validated cspell dependency.

- Reconcile standalone dependency checksums against deterministic current
  module archives so CI, local verification, and release consumers resolve
  identical content.

- Harden standalone documentation validation with deterministic spelling and
  link checks, package-specific documentation gates, and repository-local
  contributor guidance.

### Documentation

- Replace obsolete standalone-repository links and workflow claims with
  monorepo-canonical targets and current release guidance.
- Keep the initial `v1.0.0` scope under Unreleased until a tag is published.

- Link the package README to the repository-wide Golib documentation portal.

### Compatibility

- Added a pinned module export baseline so incompatible public API changes
  fail the canonical repository gate.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-webhook` identity while preserving its documented API and behavior.
- Refresh local `v0.0.0` owned-module checksums after dependency manifests and
  release notes were normalized; runtime behavior and public APIs are
  unchanged.
- Align the transitive `golang.org/x/text` dependency with v0.41.0 after the
  owned module graph removed GO-2026-5970.
- Rename the unpublished adapter packages to target-oriented `idempotency`,
  `slog`, `outbox`, `queue`, and `otel` paths with unambiguous `webhook*`
  package identifiers.
- Declare the unresolved-decision inventory explicitly so repository
  specification governance fails closed on any future open interpretation.
- Link the specification source matrix directly to the canonical decision
  register.
- Require owned sibling modules at local `v0.0.0`; clean external consumers
  pin each module to an exact main pseudo-version.

- Upgrade gRPC to 1.82.1 to remove the `GO-2026-6061` vulnerabilities and
  align the isolated module graph.
- Refresh owned-module checksums against the final consolidated archives.
- Normalized standalone module metadata against the canonical owned dependency
  graph, including complete checksums for clean consumer resolution.
- Refreshed the canonical logging checksum after its API compatibility tooling
  was standardized.

### Added

- Added an auditable specification decision register, integrity-pinned
  normative source manifest, and executable conformance gate for signature,
  HTTP, replay, envelope, delivery, retry, and SSRF policies.
- Versioned HMAC-SHA-256 and HMAC-SHA-512 signing and verification.
- Signed, bounded nonces with injectable generation and a `crypto/rand` default.
- Exact-byte HTTP verification, bounded bodies and headers, rotation windows,
  safe typed failures, middleware, atomic replay protection, and a
  `idempotency` adapter.
- Deterministic envelopes, bounded delivery and retries, `Retry-After`,
  dead-letter and replay hooks, fan-out, SSRF and DNS-rebinding protection,
  and `queue` and `outbox` adapters.
- Secret-safe observations, independent Python vectors, fuzzers, allocation
  benchmarks, complete production coverage, and release gates.
- Compiled `log` diagnostics, telemetry HTTP propagation, deterministic
  consumer fixtures, and an executable queued-delivery example.
- Pinned GitHub Actions workflow linting in local and CI release gates.
- Enforced a pure-Go dependency graph in the standalone safety gate.

### Fixed

- Reject malformed and nonnumeric endpoint ports before DNS resolution or
  dialing.
- Isolate each verification candidate's timestamp and nonce so an invalid
  rotation signature cannot alter a later valid candidate.
- Reject duplicate signing key IDs and independently validate delivery,
  wire, replay, queue, outbox, logging, and telemetry boundaries.
- Check response-body and telemetry-runtime cleanup failures in adapter and
  SSRF tests instead of silently discarding them.
- Express the accepted HTTPS or explicitly enabled HTTP schemes directly,
  preserving the default-deny SSRF policy without negation ambiguity.
- Saturate numeric `Retry-After` values at the exact configured maximum,
  including subsecond limits, and reject body limits whose sentinel byte would
  overflow an `int64` bound.
- Compare caller timestamps at the signature protocol's Unix-second precision.
- Select rotation keys at the signed timestamp, reject negative timestamps,
  and reject inverted key validity windows.
- Keep replay identity stable across overlapping secret rotation keys.
- Cover bounded `Content-Type` and `Idempotency-Key` values in v1 signatures.
- Preserve and authenticate duplicate query-value order while sorting keys.
- Preserve and authenticate the exact case-sensitive HTTP method.
- Clamp delivery latency observations when an injected clock moves backward.

### Planned v1.0.0 scope

The first release will freeze the `v1` canonicalization and wire contracts.

[Unreleased]: https://github.com/faustbrian/go-webhook/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/faustbrian/go-webhook/releases/tag/v1.0.0
