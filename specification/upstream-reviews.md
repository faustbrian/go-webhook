# Upstream authority review history

This append-only record preserves reviewed changes to the authorities monitored
by [`monitoring.json`](monitoring.json). A monitoring digest changes only after
the corresponding upstream delta has been classified against the applicable
specification decisions.

## 2026-09-28: RFC 9110 errata

- **Authority:** `rfc9110-errata`
- **URL:** https://errata.rfc-editor.org/search/?rfc_number=9110&presentation=records
- **Previous SHA-256:**
  `1f6790054c0cdb2f2a70a94fa2b9c73b09a4ee0578a32b4a3006ed0ecfaac86d`
- **Reviewed SHA-256:**
  `cec32fd170146656d933f627b512f2e027ae5c3592f5ec7760c3627493b30505`
- **Retrieved and reviewed:** 2026-09-28
- **Applicability:** `WEBHOOK-DEC-005`, `WEBHOOK-DEC-012`, and
  `WEBHOOK-DEC-015` directly; `WEBHOOK-DEC-014` as outbound header
  ownership context
- **Disposition:** Behavior-neutral. The reported correction describes
  equivalent ABNF normalization in RFC 9110 Appendix A; Webhook neither parses
  that grammar nor serializes combined repeated header fields.

[Errata ID 9164](https://errata.rfc-editor.org/eid9164/) was reported on
2026-09-07. It documents that Appendix A combines adjacent quoted terminals,
uses hexadecimal notation for case-sensitive strings, and shortens equivalent
repetition bounds in addition to the declared list-rule expansion. The report
states that the collected grammar and body define the same language. It is
still Reported, not Verified, and the immutable RFC 9110 source is unchanged.

Webhook requires singleton authentication and event-ID headers, rejects
duplicates instead of combining them, and delegates outbound HTTP wire
serialization to `net/http`. For `WEBHOOK-DEC-015`, the equivalent `GMT` and
`date3` ABNF spellings do not change the selected `Retry-After` HTTP-date
behavior, which delegates parsing to `http.ParseTime`. The existing
`WEBHOOK-DEC-005`, `-012`, `-014`, and `-015` decisions and their conformance
bindings therefore remain selected.
Errata ID 9162 is also still Reported. Reconsider if a relevant erratum is
Verified with a semantic change or Webhook takes ownership of HTTP grammar or
combined-field serialization.

## 2026-09-03: RFC 9110 errata

- **Authority:** `rfc9110-errata`
- **URL:** https://errata.rfc-editor.org/search/?rfc_number=9110&presentation=records
- **Previous SHA-256:**
  `38bd006c96f8963d58573f704c5313a5f81968b90738c03ade0b036ec7bbdf4b`
- **Reviewed SHA-256:**
  `1f6790054c0cdb2f2a70a94fa2b9c73b09a4ee0578a32b4a3006ed0ecfaac86d`
- **Retrieved and reviewed:** 2026-09-03
- **Applicability:** `WEBHOOK-DEC-005` and `WEBHOOK-DEC-012` directly;
  `WEBHOOK-DEC-014` as outbound header ownership context
- **Disposition:** Behavior-neutral because webhook authentication rejects
  repeated singleton fields and does not depend on combined-field
  serialization.

[Errata ID 9162](https://errata.rfc-editor.org/eid9162/) was reported on
2026-09-01 as a Technical erratum against RFC 9110 Section 5.2. It proposes
changing the repeated-field combination wording from values separated by a
comma to values separated by comma plus space so the rule matches its example.
The erratum remains Reported, not Verified, and does not revise the immutable
RFC 9110 source.

`WEBHOOK-DEC-005` requires exactly one `Content-Type` and `Idempotency-Key`
value, and `WEBHOOK-DEC-012` requires exactly one event-ID field. Duplicate
values are rejected instead of combined. `WEBHOOK-DEC-014` sets the two owned
outbound singleton fields and otherwise leaves cloned caller headers to
`net/http`; no selected behavior or source binding changes. Reconsider this
disposition if Errata ID 9162 becomes Verified or if the package later owns
serialization of combined repeated fields.
