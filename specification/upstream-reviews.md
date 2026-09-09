# Upstream authority review history

This append-only record preserves reviewed changes to the authorities monitored
by [`monitoring.json`](monitoring.json). A monitoring digest changes only after
the corresponding upstream delta has been classified against the applicable
specification decisions.

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

## 2026-09-09: RFC 9110 errata

- **Authority:** `rfc9110-errata`
- **URL:** https://errata.rfc-editor.org/search/?rfc_number=9110&presentation=records
- **Previous SHA-256:**
  `1f6790054c0cdb2f2a70a94fa2b9c73b09a4ee0578a32b4a3006ed0ecfaac86d`
- **Reviewed SHA-256:**
  `cec32fd170146656d933f627b512f2e027ae5c3592f5ec7760c3627493b30505`
- **Retrieved and reviewed:** 2026-09-09
- **Applicability:** No selected webhook behavior; the package does not derive
  or mechanically compare the collected ABNF in RFC 9110 Appendix A.
- **Disposition:** Behavior-neutral clarification of equivalent ABNF
  renderings.

[Errata ID 9164](https://errata.rfc-editor.org/eid9164/) was reported on
2026-09-07 as a Technical erratum against RFC 9110 Appendix A. It proposes
documenting that the collected ABNF normalizes adjacent terminals,
case-sensitive strings, and repetition bounds in addition to expanding list
rules. The erratum states that these renderings define the same language and
remains Reported, not Verified.

Webhook decisions cite behavioral RFC 9110 sections and define their own
bounded header, request, and response contracts. They neither consume the
collected ABNF mechanically nor depend on the spelling of the equivalent
grammar forms, so no source, test, or selected behavior changes. Reconsider
this disposition if the erratum changes the language recognized by RFC 9110
or webhook begins generating a parser from Appendix A.

## 2026-09-09: W3C Trace Context publication history

- **Authority:** `trace-context-releases`
- **URL:** https://www.w3.org/standards/history/trace-context/
- **Previous SHA-256:**
  `210f1d71d5a667a18beea7c681e484088023ed8d839106f757b411c644cdfeaf`
- **Reviewed SHA-256:**
  `ab1a5ce3b8491194038e5b0b7c5a1b06bd2b1660d29b50caf841a46f7d687968`
- **Retrieved and reviewed:** 2026-09-09
- **Applicability:** `WEBHOOK-DEC-023` retains its immutable W3C Trace Context
  Level 1 Recommendation binding.
- **Disposition:** Behavior-neutral page-chrome change; the publication list
  and selected immutable Recommendation are unchanged.

The history page still identifies the 23 November 2021 Recommendation as its
latest publication and lists the same seven Trace Context publications. The
selected Recommendation bytes retain their pinned SHA-256
`9e7228d2a91c5aa4bef6e7f610366a2000274ee51699053e10c8ac3f0b8965be`.
The current page adds mutable W3C site chrome, including a 2026 community
survey banner and refreshed copyright year; neither changes the standard or
webhook trace propagation. Reconsider this disposition when the publication
list or selected Recommendation changes.
