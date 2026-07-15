# Hardening findings

| Severity | Finding | Evidence and disposition |
| --- | --- | --- |
| Critical | Redirect and DNS changes could bypass URL-only policy | Resolved: redirects refused; all DNS answers validated again in the dialer |
| High | Concurrent replay acceptance | Resolved: atomic store contract and 64-way race test accepts exactly one |
| High | Ambiguous delivery could multiply effects | Resolved: retries require idempotency; queue/outbox handlers use one attempt |
| Medium | Malformed/duplicate headers could parse ambiguously | Resolved: strict bounded structured parser and mutation/fuzz tests |
| Medium | Body/response/fan-out resource exhaustion | Resolved: hard byte, attempt, address, and worker bounds with fault tests |
| Medium | Diagnostics could expose attacker-controlled secrets | Resolved: safe typed errors and closed observation schema |
| Low | Named HTTP client integration unavailable | Accepted: repository has no published module/API; narrow `HTTPDoer` seam documented |

No open critical, high, or medium finding remains. Provider support remains
intentionally empty rather than making a claim without authoritative vectors.
