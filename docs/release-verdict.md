# Release and verification record

## Published result

`v1.0.0` was published on 2026-08-26. The published
[GitHub release](https://github.com/faustbrian/go-webhook/releases/tag/v1.0.0)
records the released module version. Its final candidate passed the complete
repository gate, module verification, production statement coverage, race and
fuzz targets, the independent Python fixture, workflow lint, and vulnerability
scanning.

No critical, high, or medium finding remained open at publication. The generic
SHA-256 and SHA-512 schemes had independent vectors, and the provider matrix
intentionally claimed no provider preset without authoritative conformance
evidence.

## Release verification commands

Release evidence is collected from a clean final tree with:

```sh
make check FUZZTIME=10s
go mod verify
git status --porcelain=v1
```

`make check` aggregates format, vet, static analysis, tests, meaningful 100%
production coverage, race, bounded fuzzing, allocation benchmarks, executable
documentation, `GO-SAFETY-1`, independent interoperability, and vulnerability
scanning. GitHub Actions runs the repository and module gates.

## Verdict criteria

The verdict is `GO` only when all commands exit zero, the worktree is clean,
the changelog moves the planned scope under the release version, no
critical/high/medium finding is open, and provider claims exactly match the
provider matrix. Otherwise it is `NO-GO`.

Residual risks are operator-owned secret quality and rotation, durable store
availability and tenant scoping, application-side idempotency and payload
validation, and explicit weakening through private SSRF allow prefixes. The
module makes no exactly-once claim.
