# Webhook Specification Decisions

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and
"OPTIONAL" in this document are to be interpreted as described in BCP 14
[RFC2119] [RFC8174] when, and only when, they appear in all capitals, as
shown here.

This register records the observable choices made by the generic `v1`
signature protocol and its HTTP, replay, delivery, envelope, and endpoint
security boundaries. `webhook` defines a local protocol. It does not claim
compatibility with RFC 9421, CloudEvents, an Idempotency-Key draft, or any
vendor webhook scheme unless a future isolated adapter says so explicitly.

Every resolved entry names executable evidence. A change to a selected
behavior requires compatibility, security, resource, conformance, API, and
changelog review. Superseded entries remain in this file and link to their
replacement.

## WEBHOOK-DEC-001: Protocol identity and HTTP Message Signatures

**Authoritative reference:** [RFC 9421](https://www.rfc-editor.org/rfc/rfc9421.html).

- **Status, owner, and classification:** `resolved`; `webhook` maintainers;
  local wire protocol and compatibility policy.
- **Source and issue:** RFC 9421 [HTTP Message Signatures](https://www.rfc-editor.org/rfc/rfc9421.html)
  defines signature parameters, covered components, derived components, and
  Structured Fields serialization. Existing webhook providers instead use
  mutually incompatible proprietary canonical forms. The package must not
  imply that a custom webhook MAC is RFC 9421 or vendor compatible.
- **Interpretations and peer behavior:** Implement RFC 9421 directly, clone a
  provider format, expose an unversioned helper, or define a small versioned
  local protocol. Providers disagree on headers, timestamps, body treatment,
  key selection, and canonical request targets.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  `Webhook-Signature` `v1` is a local,
  versioned HMAC protocol specified by `docs/signatures.md`. It neither parses
  nor emits RFC 9421 `Signature` or `Signature-Input` fields and makes no
  provider claim. Negotiation is not implicit: unknown versions fail closed.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestCanonicalizeProducesStableVersionedBytes`,
  `TestParseSignatureHeadersRejectsMalformedOrAmbiguousInput`, and
  `TestIndependentInteroperabilityVectors` cover `Canonicalize`, `Signature`,
  and the HTTP helpers. There is no upstream issue because this is deliberate
  protocol ownership. Reconsider only through a separately versioned RFC 9421
  or provider adapter with independent vectors.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-001","title":"Protocol identity and HTTP Message Signatures","status":"resolved","owner":"`webhook` maintainers","classification":"interoperability policy","decision_scope":"application-policy","specification":"RFC 9421 HTTP Message Signatures comparison boundary","version":"RFC 9421","source_authority":"rfc9421-source","section":"RFC 9421 Sections 2-4","requirement_strength":"not specified","issue":"RFC 9421 [HTTP Message Signatures](https://www.rfc-editor.org/rfc/rfc9421.html) defines signature parameters, covered components, derived components, and Structured Fields serialization. Existing webhook providers instead use mutually incompatible proprietary canonical forms. The package must not imply that a custom webhook MAC is RFC 9421 or vendor compatible.","interpretations":["Implement RFC 9421 directly, clone a provider format, expose an unversioned helper, or define a small versioned local protocol. Providers disagree on headers, timestamps, body treatment, key selection, and canonical request targets."],"peer_behavior":"Maintained-provider behavior has not been assessed for this decision; the independent Python standard-library generator agrees with the pinned local v1 vectors.","selected_behavior":"`Webhook-Signature` `v1` is a local, versioned HMAC protocol specified by `docs/signatures.md`. It neither parses nor emits RFC 9421 `Signature` or `Signature-Input` fields and makes no provider claim. Negotiation is not implicit: unknown versions fail closed.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestCanonicalizeProducesStableVersionedBytes","TestParseSignatureHeadersRejectsMalformedOrAmbiguousInput","TestIndependentInteroperabilityVectors"],"fixture_evidence":["testdata/vectors/v1.json"],"fuzz_evidence":["FuzzCanonicalize","FuzzParseSignatureHeaders"],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["Canonicalize","Signature","SignRequest","VerifyRequest"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc9421.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
interoperability policy
application-policy
RFC 9421 HTTP Message Signatures comparison boundary
RFC 9421
rfc9421-source
https://www.rfc-editor.org/rfc/rfc9421.txt
RFC 9421 Sections 2-4
not specified
RFC 9421 [HTTP Message Signatures](https://www.rfc-editor.org/rfc/rfc9421.html) defines signature parameters, covered components, derived components, and Structured Fields serialization. Existing webhook providers instead use mutually incompatible proprietary canonical forms. The package must not imply that a custom webhook MAC is RFC 9421 or vendor compatible.
Maintained-provider behavior has not been assessed for this decision; the independent Python standard-library generator agrees with the pinned local v1 vectors.
`Webhook-Signature` `v1` is a local, versioned HMAC protocol specified by `docs/signatures.md`. It neither parses nor emits RFC 9421 `Signature` or `Signature-Input` fields and makes no provider claim. Negotiation is not implicit: unknown versions fail closed.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Implement RFC 9421 directly, clone a provider format, expose an unversioned helper, or define a small versioned local protocol. Providers disagree on headers, timestamps, body treatment, key selection, and canonical request targets.
TestCanonicalizeProducesStableVersionedBytes
TestParseSignatureHeadersRejectsMalformedOrAmbiguousInput
TestIndependentInteroperabilityVectors
testdata/vectors/v1.json
FuzzCanonicalize
FuzzParseSignatureHeaders
Canonicalize
Signature
SignRequest
VerifyRequest
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-002: HMAC algorithms and body digest

**Authoritative reference:** [RFC 2104](https://www.rfc-editor.org/rfc/rfc2104.html).

- **Status, owner, and classification:** `resolved`; maintainers; normative
  cryptographic profile plus local algorithm policy.
- **Source and issue:** RFC 2104 defines HMAC, while RFC 4231 supplies
  HMAC-SHA-256 and HMAC-SHA-512 test cases. Neither chooses a webhook algorithm,
  body digest, downgrade policy, or algorithm-negotiation mechanism.
- **Interpretations and peer behavior:** Support one digest, permit arbitrary
  hash names, derive the body digest from the MAC algorithm, or fix a body
  digest independently. Provider schemes use several incompatible choices.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  A signer and verifier are configured
  for exactly one allow-listed `sha256` or `sha512` HMAC algorithm. The canonical
  body component is always SHA-256, including under HMAC-SHA-512. The algorithm
  name is signed, unknown algorithms are rejected, and no inbound header can
  widen the configured algorithm policy.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestSignerAndVerifierSupportSHA256AndSHA512`,
  `TestVerifierRejectsMutationOfEverySignedComponent`, and the RFC-independent
  Python vectors cover `Algorithm`, `Signer`, and `Verifier`. Reconsider for a
  concrete cryptographic migration with a new protocol version and downgrade
  analysis.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-002","title":"HMAC algorithms and body digest","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"defensive","specification":"RFC 2104 HMAC and RFC 4231 HMAC-SHA-256/HMAC-SHA-512 vectors","version":"RFC 2104","source_authority":"rfc2104-source","section":"RFC 2104 Sections 2-3; RFC 4231 Sections 4-5","requirement_strength":"not specified","issue":"RFC 2104 defines HMAC, while RFC 4231 supplies HMAC-SHA-256 and HMAC-SHA-512 test cases. Neither chooses a webhook algorithm, body digest, downgrade policy, or algorithm-negotiation mechanism.","interpretations":["Support one digest, permit arbitrary hash names, derive the body digest from the MAC algorithm, or fix a body digest independently. Provider schemes use several incompatible choices."],"peer_behavior":"Maintained-provider behavior has not been assessed for this decision; the independent Python standard-library generator agrees with the pinned local v1 vectors.","selected_behavior":"A signer and verifier are configured for exactly one allow-listed `sha256` or `sha512` HMAC algorithm. The canonical body component is always SHA-256, including under HMAC-SHA-512. The algorithm name is signed, unknown algorithms are rejected, and no inbound header can widen the configured algorithm policy.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestSignerAndVerifierSupportSHA256AndSHA512","TestVerifierRejectsMutationOfEverySignedComponent","TestIndependentInteroperabilityVectors"],"fixture_evidence":["testdata/vectors/v1.json"],"fuzz_evidence":["FuzzCanonicalize"],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["Algorithm","Signer","Verifier"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc2104.txt

Additional authoritative source: `{"id":"rfc4231-source","version":"RFC 4231","url":"https://www.rfc-editor.org/rfc/rfc4231.txt","specifications":["RFC 2104 HMAC and RFC 4231 HMAC-SHA-256/HMAC-SHA-512 vectors"]}`

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
defensive
RFC 2104 HMAC and RFC 4231 HMAC-SHA-256/HMAC-SHA-512 vectors
RFC 2104
rfc2104-source
https://www.rfc-editor.org/rfc/rfc2104.txt
RFC 2104 Sections 2-3; RFC 4231 Sections 4-5
not specified
RFC 2104 defines HMAC, while RFC 4231 supplies HMAC-SHA-256 and HMAC-SHA-512 test cases. Neither chooses a webhook algorithm, body digest, downgrade policy, or algorithm-negotiation mechanism.
Maintained-provider behavior has not been assessed for this decision; the independent Python standard-library generator agrees with the pinned local v1 vectors.
A signer and verifier are configured for exactly one allow-listed `sha256` or `sha512` HMAC algorithm. The canonical body component is always SHA-256, including under HMAC-SHA-512. The algorithm name is signed, unknown algorithms are rejected, and no inbound header can widen the configured algorithm policy.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Support one digest, permit arbitrary hash names, derive the body digest from the MAC algorithm, or fix a body digest independently. Provider schemes use several incompatible choices.
TestSignerAndVerifierSupportSHA256AndSHA512
TestVerifierRejectsMutationOfEverySignedComponent
TestIndependentInteroperabilityVectors
testdata/vectors/v1.json
FuzzCanonicalize
Algorithm
Signer
Verifier
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-003: Canonical framing, encoding, and Unicode

**Authoritative reference:** [RFC 4648](https://www.rfc-editor.org/rfc/rfc4648.html).

- **Status, owner, and classification:** `resolved`; maintainers; local
  wire-format policy informed by RFC 4648.
- **Source and issue:** RFC 4648 Section 5 defines base64url and Section 3.2
  discusses padding, but no standard defines canonical webhook field framing,
  line endings, Unicode normalization, or empty-field representation.
- **Interpretations and peer behavior:** Concatenate raw values with separators,
  serialize JSON, use length prefixes, normalize Unicode, or encode every
  variable field. Ambiguous concatenation and platform line endings can create
  cross-language disagreement.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  `v1` uses fixed ASCII labels in a
  fixed order, LF delimiters, a final empty line, and unpadded base64url for
  every variable byte field. It signs caller string bytes exactly and performs
  no Unicode normalization. Empty values remain explicit empty encoded fields.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestCanonicalizeProducesStableVersionedBytes`, `FuzzCanonicalize`, and
  `TestIndependentInteroperabilityVectors` cover `Canonicalize` and the golden
  wire bytes. Reconsider only with a new version and cross-language vectors for
  every byte-level change.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-003","title":"Canonical framing, encoding, and Unicode","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"application-policy","specification":"RFC 4648 Base-N Encodings","version":"RFC 4648","source_authority":"rfc4648-source","section":"RFC 4648 Sections 3.2 and 5","requirement_strength":"not specified","issue":"RFC 4648 Section 5 defines base64url and Section 3.2 discusses padding, but no standard defines canonical webhook field framing, line endings, Unicode normalization, or empty-field representation.","interpretations":["Concatenate raw values with separators, serialize JSON, use length prefixes, normalize Unicode, or encode every variable field. Ambiguous concatenation and platform line endings can create cross-language disagreement."],"peer_behavior":"Maintained-provider behavior has not been assessed for this decision; the independent Python standard-library generator agrees with the pinned local v1 vectors.","selected_behavior":"`v1` uses fixed ASCII labels in a fixed order, LF delimiters, a final empty line, and unpadded base64url for every variable byte field. It signs caller string bytes exactly and performs no Unicode normalization. Empty values remain explicit empty encoded fields.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestCanonicalizeProducesStableVersionedBytes","TestIndependentInteroperabilityVectors"],"fixture_evidence":["testdata/vectors/v1.json"],"fuzz_evidence":["FuzzCanonicalize"],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["Canonicalize","Signature"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc4648.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
application-policy
RFC 4648 Base-N Encodings
RFC 4648
rfc4648-source
https://www.rfc-editor.org/rfc/rfc4648.txt
RFC 4648 Sections 3.2 and 5
not specified
RFC 4648 Section 5 defines base64url and Section 3.2 discusses padding, but no standard defines canonical webhook field framing, line endings, Unicode normalization, or empty-field representation.
Maintained-provider behavior has not been assessed for this decision; the independent Python standard-library generator agrees with the pinned local v1 vectors.
`v1` uses fixed ASCII labels in a fixed order, LF delimiters, a final empty line, and unpadded base64url for every variable byte field. It signs caller string bytes exactly and performs no Unicode normalization. Empty values remain explicit empty encoded fields.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Concatenate raw values with separators, serialize JSON, use length prefixes, normalize Unicode, or encode every variable field. Ambiguous concatenation and platform line endings can create cross-language disagreement.
TestCanonicalizeProducesStableVersionedBytes
TestIndependentInteroperabilityVectors
testdata/vectors/v1.json
FuzzCanonicalize
Canonicalize
Signature
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-004: Method, path, and query canonicalization

**Authoritative reference:** [RFC 3986](https://www.rfc-editor.org/rfc/rfc3986.html).

- **Status, owner, and classification:** `resolved`; maintainers; Go URL
  interoperability and local request-target policy.
- **Source and issue:** RFC 3986 defines URI components, while Go 1.26.6
  `net/url` defines `EscapedPath`, form-style query parsing, and `Values.Encode`.
  Neither selects canonical webhook semantics for method case, percent escapes,
  plus signs, duplicate values, or an empty path.
- **Interpretations and peer behavior:** Sign the raw request target, decoded
  values, a reconstructed URL, an uppercased method, or a provider-specific
  subset. Peers differ especially around duplicate query parameters.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  Sign the exact nonempty method with
  case preserved; use `EscapedPath` with `/` for an empty path; parse the raw
  query using Go form semantics; sort keys through `Values.Encode`; preserve
  duplicate value order; and emit Go's canonical escaping. Semantically
  equivalent raw query spellings can therefore share canonical bytes, while
  reordered duplicate values do not.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestVerifierBindsDuplicateQueryValueOrder`,
  `TestVerifierRejectsMutationOfEverySignedComponent`,
  `TestSignAndVerifyRequestUsesRawBodyAndRestoresIt`, and the Python vectors
  cover `Message`, `Canonicalize`, `SignRequest`, and `VerifyRequest`. Reconsider
  if Go changes these URL contracts or a raw-target protocol version is added.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-004","title":"Method, path, and query canonicalization","status":"resolved","owner":"`webhook` maintainers","classification":"ambiguity","decision_scope":"application-policy","specification":"RFC 3986 URI Generic Syntax","version":"RFC 3986","source_authority":"rfc3986-source","section":"RFC 3986 Sections 2-5","requirement_strength":"not specified","issue":"RFC 3986 defines URI components, while Go 1.26.6 `net/url` defines `EscapedPath`, form-style query parsing, and `Values.Encode`. Neither selects canonical webhook semantics for method case, percent escapes, plus signs, duplicate values, or an empty path.","interpretations":["Sign the raw request target, decoded values, a reconstructed URL, an uppercased method, or a provider-specific subset. Peers differ especially around duplicate query parameters."],"peer_behavior":"Maintained-provider behavior has not been assessed for this decision; the independent Python standard-library generator agrees with the pinned local v1 vectors.","selected_behavior":"Sign the exact nonempty method with case preserved; use `EscapedPath` with `/` for an empty path; parse the raw query using Go form semantics; sort keys through `Values.Encode`; preserve duplicate value order; and emit Go's canonical escaping. Semantically equivalent raw query spellings can therefore share canonical bytes, while reordered duplicate values do not.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestVerifierBindsDuplicateQueryValueOrder","TestVerifierRejectsMutationOfEverySignedComponent","TestSignAndVerifyRequestUsesRawBodyAndRestoresIt","TestIndependentInteroperabilityVectors"],"fixture_evidence":["testdata/vectors/v1.json"],"fuzz_evidence":["FuzzCanonicalize"],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["Message","Canonicalize","SignRequest","VerifyRequest"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc3986.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
ambiguity
application-policy
RFC 3986 URI Generic Syntax
RFC 3986
rfc3986-source
https://www.rfc-editor.org/rfc/rfc3986.txt
RFC 3986 Sections 2-5
not specified
RFC 3986 defines URI components, while Go 1.26.6 `net/url` defines `EscapedPath`, form-style query parsing, and `Values.Encode`. Neither selects canonical webhook semantics for method case, percent escapes, plus signs, duplicate values, or an empty path.
Maintained-provider behavior has not been assessed for this decision; the independent Python standard-library generator agrees with the pinned local v1 vectors.
Sign the exact nonempty method with case preserved; use `EscapedPath` with `/` for an empty path; parse the raw query using Go form semantics; sort keys through `Values.Encode`; preserve duplicate value order; and emit Go's canonical escaping. Semantically equivalent raw query spellings can therefore share canonical bytes, while reordered duplicate values do not.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Sign the raw request target, decoded values, a reconstructed URL, an uppercased method, or a provider-specific subset. Peers differ especially around duplicate query parameters.
TestVerifierBindsDuplicateQueryValueOrder
TestVerifierRejectsMutationOfEverySignedComponent
TestSignAndVerifyRequestUsesRawBodyAndRestoresIt
TestIndependentInteroperabilityVectors
testdata/vectors/v1.json
FuzzCanonicalize
Message
Canonicalize
SignRequest
VerifyRequest
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-005: Host and behavior-changing header coverage

**Authoritative reference:** [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110.html).

- **Status, owner, and classification:** `resolved`; maintainers; defensive
  HTTP field and canonicalization policy.
- **Source and issue:** RFC 9110 defines the target URI, media type field, and
  field combination rules but does not define a webhook canonical host or the
  application-defined `Idempotency-Key` field used here.
- **Interpretations and peer behavior:** Exclude host and application fields,
  lowercase the full authority, remove default ports, combine duplicate field
  lines, or bind exact values. Intermediaries can otherwise alter decoding or
  deduplication semantics without changing the body.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  Lowercase the request host while
  preserving an explicit port. Bind exactly one UTF-8 `Content-Type` and
  `Idempotency-Key` value, each bounded to 256 bytes; absence is an explicit
  empty canonical field. Duplicate, line-breaking, oversized, or invalid UTF-8
  values fail before body capture. Trace fields remain outside the signature.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestVerifyRequestRejectsDuplicateSignedHeaderBeforeBodyRead`,
  `TestVerifyRequestRejectsMutationOfFixedSignedHeaders`, and
  `TestFixedHeadersAndEventIDsPreserveExactBoundaries` cover `RequestOptions`,
  `SignRequest`, and `VerifyRequest`. Reconsider for a new signed-field profile
  or if the package adopts a standardized idempotency field contract.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-005","title":"Host and behavior-changing header coverage","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"defensive","specification":"RFC 9110 HTTP Semantics","version":"RFC 9110","source_authority":"rfc9110-source","section":"RFC 9110 Sections 5-7","requirement_strength":"not specified","issue":"RFC 9110 defines the target URI, media type field, and field combination rules but does not define a webhook canonical host or the application-defined `Idempotency-Key` field used here.","interpretations":["Exclude host and application fields, lowercase the full authority, remove default ports, combine duplicate field lines, or bind exact values. Intermediaries can otherwise alter decoding or deduplication semantics without changing the body."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"Lowercase the request host while preserving an explicit port. Bind exactly one UTF-8 `Content-Type` and `Idempotency-Key` value, each bounded to 256 bytes; absence is an explicit empty canonical field. Duplicate, line-breaking, oversized, or invalid UTF-8 values fail before body capture. Trace fields remain outside the signature.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestVerifyRequestRejectsMutationOfFixedSignedHeaders","TestVerifyRequestRejectsDuplicateSignedHeaderBeforeBodyRead"],"fixture_evidence":[],"fuzz_evidence":["FuzzParseSignatureHeaders"],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["Message","SignRequest","VerifyRequest"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc9110.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
defensive
RFC 9110 HTTP Semantics
RFC 9110
rfc9110-source
https://www.rfc-editor.org/rfc/rfc9110.txt
RFC 9110 Sections 5-7
not specified
RFC 9110 defines the target URI, media type field, and field combination rules but does not define a webhook canonical host or the application-defined `Idempotency-Key` field used here.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
Lowercase the request host while preserving an explicit port. Bind exactly one UTF-8 `Content-Type` and `Idempotency-Key` value, each bounded to 256 bytes; absence is an explicit empty canonical field. Duplicate, line-breaking, oversized, or invalid UTF-8 values fail before body capture. Trace fields remain outside the signature.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Exclude host and application fields, lowercase the full authority, remove default ports, combine duplicate field lines, or bind exact values. Intermediaries can otherwise alter decoding or deduplication semantics without changing the body.
TestVerifyRequestRejectsMutationOfFixedSignedHeaders
TestVerifyRequestRejectsDuplicateSignedHeaderBeforeBodyRead
FuzzParseSignatureHeaders
Message
SignRequest
VerifyRequest
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-006: Timestamp precision and tolerance

**Authoritative reference:** [RFC 3339](https://www.rfc-editor.org/rfc/rfc3339.html).

- **Status, owner, and classification:** `resolved`; maintainers; local replay
  freshness policy using Go time semantics.
- **Source and issue:** HTTP dates in RFC 9110 and Internet timestamps in RFC
  3339 do not define a webhook signature timestamp, allowed skew, precision, or
  whether tolerance endpoints are inclusive.
- **Interpretations and peer behavior:** Use milliseconds, seconds, HTTP-date,
  asymmetric age limits, or exact equality. Providers vary and often leave
  boundary behavior undocumented.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  `v1` signs a canonical nonnegative
  Unix-second integer. Signer inputs are truncated to that precision. Verifiers
  accept absolute skew less than or equal to the configured tolerance, compare
  caller-provided timestamps by Unix second, and reject negative,
  noncanonical, overflowing, or outside-window values.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestVerifierTimestampToleranceBoundaries`,
  `TestVerifierComparesCallerTimestampAtProtocolSecondPrecision`,
  `FuzzTimestampVerification`, and header parser tests cover `Signature`,
  `Signer`, and `Verifier`. Reconsider only with a new protocol version or a
  demonstrated need for asymmetric age policy.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-006","title":"Timestamp precision and tolerance","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"defensive","specification":"RFC 3339 Internet date and time","version":"RFC 3339","source_authority":"rfc3339-source","section":"RFC 3339 Section 5","requirement_strength":"not specified","issue":"HTTP dates in RFC 9110 and Internet timestamps in RFC 3339 do not define a webhook signature timestamp, allowed skew, precision, or whether tolerance endpoints are inclusive.","interpretations":["Use milliseconds, seconds, HTTP-date, asymmetric age limits, or exact equality. Providers vary and often leave boundary behavior undocumented."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"`v1` signs a canonical nonnegative Unix-second integer. Signer inputs are truncated to that precision. Verifiers accept absolute skew less than or equal to the configured tolerance, compare caller-provided timestamps by Unix second, and reject negative, noncanonical, overflowing, or outside-window values.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestVerifierComparesCallerTimestampAtProtocolSecondPrecision","TestVerifierTimestampToleranceBoundaries"],"fixture_evidence":[],"fuzz_evidence":["FuzzTimestampVerification"],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["Signer","Verifier","Signature"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc3339.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
defensive
RFC 3339 Internet date and time
RFC 3339
rfc3339-source
https://www.rfc-editor.org/rfc/rfc3339.txt
RFC 3339 Section 5
not specified
HTTP dates in RFC 9110 and Internet timestamps in RFC 3339 do not define a webhook signature timestamp, allowed skew, precision, or whether tolerance endpoints are inclusive.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
`v1` signs a canonical nonnegative Unix-second integer. Signer inputs are truncated to that precision. Verifiers accept absolute skew less than or equal to the configured tolerance, compare caller-provided timestamps by Unix second, and reject negative, noncanonical, overflowing, or outside-window values.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Use milliseconds, seconds, HTTP-date, asymmetric age limits, or exact equality. Providers vary and often leave boundary behavior undocumented.
TestVerifierComparesCallerTimestampAtProtocolSecondPrecision
TestVerifierTimestampToleranceBoundaries
FuzzTimestampVerification
Signer
Verifier
Signature
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-007: Nonces, key rotation, and revocation

**Authoritative reference:** [RFC 2104](https://www.rfc-editor.org/rfc/rfc2104.html).

- **Status, owner, and classification:** `resolved`; maintainers; defensive
  cryptographic lifecycle policy.
- **Source and issue:** HMAC standards do not require a nonce, key identifier,
  rotation overlap, validity interval, ordering, or revocation behavior for
  webhook messages.
- **Interpretations and peer behavior:** Emit deterministic MACs, sign with one
  current key, emit every active key, select keys at wall-clock or signed time,
  and accept revoked keys until expiry. Provider rotation behavior varies.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  Every signature carries a nonempty
  valid UTF-8 nonce bounded to 128 bytes. The default is 18 random bytes encoded
  as base64url. One nonce is shared across all signatures in one operation.
  Signers emit every non-revoked key active at the normalized signed timestamp,
  newest first with key-ID tie-breaking; verifiers use that same timestamp.
  Duplicate IDs and inverted validity windows are configuration errors.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestSignerUsesOneInjectedNonceAcrossRotationSignatures`,
  `TestRotationSignsAllActiveKeysAndAcceptsOverlap`,
  `TestSignerSelectsRotationKeyAtSignedTimestamp`, and
  `TestSignerOrdersKeysWithEqualActivationByID` cover `SigningKey`,
  `VerificationKey`, and configuration. Reconsider for an external key service
  contract or a separately versioned single-signature negotiation scheme.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-007","title":"Nonces, key rotation, and revocation","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"defensive","specification":"RFC 2104 HMAC and RFC 4231 HMAC-SHA-256/HMAC-SHA-512 vectors","version":"RFC 2104","source_authority":"rfc2104-source","section":"RFC 2104 Sections 2-3","requirement_strength":"not specified","issue":"HMAC standards do not require a nonce, key identifier, rotation overlap, validity interval, ordering, or revocation behavior for webhook messages.","interpretations":["Emit deterministic MACs, sign with one current key, emit every active key, select keys at wall-clock or signed time, and accept revoked keys until expiry. Provider rotation behavior varies."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"Every signature carries a nonempty valid UTF-8 nonce bounded to 128 bytes. The default is 18 random bytes encoded as base64url. One nonce is shared across all signatures in one operation. Signers emit every non-revoked key active at the normalized signed timestamp, newest first with key-ID tie-breaking; verifiers use that same timestamp. Duplicate IDs and inverted validity windows are configuration errors.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestSignerUsesOneInjectedNonceAcrossRotationSignatures","TestSignerSelectsRotationKeyAtSignedTimestamp","TestRotationSignsAllActiveKeysAndAcceptsOverlap"],"fixture_evidence":[],"fuzz_evidence":["FuzzTimestampVerification"],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["Signer","Verifier","SigningKey","VerificationKey"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc2104.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
defensive
RFC 2104 HMAC and RFC 4231 HMAC-SHA-256/HMAC-SHA-512 vectors
RFC 2104
rfc2104-source
https://www.rfc-editor.org/rfc/rfc2104.txt
RFC 2104 Sections 2-3
not specified
HMAC standards do not require a nonce, key identifier, rotation overlap, validity interval, ordering, or revocation behavior for webhook messages.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
Every signature carries a nonempty valid UTF-8 nonce bounded to 128 bytes. The default is 18 random bytes encoded as base64url. One nonce is shared across all signatures in one operation. Signers emit every non-revoked key active at the normalized signed timestamp, newest first with key-ID tie-breaking; verifiers use that same timestamp. Duplicate IDs and inverted validity windows are configuration errors.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Emit deterministic MACs, sign with one current key, emit every active key, select keys at wall-clock or signed time, and accept revoked keys until expiry. Provider rotation behavior varies.
TestSignerUsesOneInjectedNonceAcrossRotationSignatures
TestSignerSelectsRotationKeyAtSignedTimestamp
TestRotationSignsAllActiveKeysAndAcceptsOverlap
FuzzTimestampVerification
Signer
Verifier
SigningKey
VerificationKey
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-008: Signature header grammar and set rejection

**Authoritative reference:** [RFC 8941](https://www.rfc-editor.org/rfc/rfc8941.html).

- **Status, owner, and classification:** `resolved`; maintainers; local HTTP
  field grammar and defensive parsing policy.
- **Source and issue:** RFC 9110 permits repeated field lines and field-specific
  combination rules. RFC 9421 uses
  [RFC 8941 Structured Fields](https://www.rfc-editor.org/rfc/rfc8941.html),
  but local `v1` does not. A custom field must define order, duplicates,
  whitespace, unknown parameters, padding, and malformed sibling behavior.
- **Interpretations and peer behavior:** Accept parameters in any order, ignore
  unknowns, split comma-combined fields, salvage valid siblings, or require
  exact serialization. Lenient parsers create cross-implementation ambiguity.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  Each active key produces one separate
  `Webhook-Signature` field line with fixed parameter order and exact lowercase
  names. Comma combination, padding, whitespace variation, duplicate or unknown
  parameters, duplicate key IDs, noncanonical timestamps, invalid encodings,
  and any malformed sibling reject the entire bounded set.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestSignatureHeadersRoundTripMultipleRotationSignatures`,
  `TestParseSignatureHeadersRejectsMalformedOrAmbiguousInput`,
  `TestParseSignatureHeadersAppliesLimitsBeforeDecoding`, and
  `FuzzParseSignatureHeaders` cover `SetSignatureHeaders` and
  `ParseSignatureHeaders`. Reconsider only in a new version with independent
  parser differential evidence.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-008","title":"Signature header grammar and set rejection","status":"resolved","owner":"`webhook` maintainers","classification":"interoperability policy","decision_scope":"defensive","specification":"RFC 8941 Structured Fields comparison boundary","version":"RFC 8941","source_authority":"rfc8941-source","section":"RFC 8941 Sections 2-4","requirement_strength":"not specified","issue":"RFC 9110 permits repeated field lines and field-specific combination rules. RFC 9421 uses [RFC 8941 Structured Fields](https://www.rfc-editor.org/rfc/rfc8941.html), but local `v1` does not. A custom field must define order, duplicates, whitespace, unknown parameters, padding, and malformed sibling behavior.","interpretations":["Accept parameters in any order, ignore unknowns, split comma-combined fields, salvage valid siblings, or require exact serialization. Lenient parsers create cross-implementation ambiguity."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"Each active key produces one separate `Webhook-Signature` field line with fixed parameter order and exact lowercase names. Comma combination, padding, whitespace variation, duplicate or unknown parameters, duplicate key IDs, noncanonical timestamps, invalid encodings, and any malformed sibling reject the entire bounded set.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestParseSignatureHeadersRejectsMalformedOrAmbiguousInput","TestParseSignatureHeadersAppliesLimitsBeforeDecoding","TestSignatureHeadersRoundTripMultipleRotationSignatures"],"fixture_evidence":[],"fuzz_evidence":["FuzzParseSignatureHeaders"],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["ParseSignatureHeaders","SignatureHeaders"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc8941.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
interoperability policy
defensive
RFC 8941 Structured Fields comparison boundary
RFC 8941
rfc8941-source
https://www.rfc-editor.org/rfc/rfc8941.txt
RFC 8941 Sections 2-4
not specified
RFC 9110 permits repeated field lines and field-specific combination rules. RFC 9421 uses [RFC 8941 Structured Fields](https://www.rfc-editor.org/rfc/rfc8941.html), but local `v1` does not. A custom field must define order, duplicates, whitespace, unknown parameters, padding, and malformed sibling behavior.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
Each active key produces one separate `Webhook-Signature` field line with fixed parameter order and exact lowercase names. Comma combination, padding, whitespace variation, duplicate or unknown parameters, duplicate key IDs, noncanonical timestamps, invalid encodings, and any malformed sibling reject the entire bounded set.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Accept parameters in any order, ignore unknowns, split comma-combined fields, salvage valid siblings, or require exact serialization. Lenient parsers create cross-implementation ambiguity.
TestParseSignatureHeadersRejectsMalformedOrAmbiguousInput
TestParseSignatureHeadersAppliesLimitsBeforeDecoding
TestSignatureHeadersRoundTripMultipleRotationSignatures
FuzzParseSignatureHeaders
ParseSignatureHeaders
SignatureHeaders
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-009: Exact request body, ordering, and ownership

**Authoritative reference:** [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110.html).

- **Status, owner, and classification:** `resolved`; maintainers; Go HTTP body
  lifecycle plus defensive resource policy.
- **Source and issue:** RFC 9110 defines message content and Go `net/http`
  exposes a stream. Neither recovers bytes consumed by earlier middleware nor
  decides whether compressed content is decoded before authentication.
- **Interpretations and peer behavior:** Hash decoded JSON, hash decompressed
  content, trust `Content-Length`, buffer without a cap, or hash exact bytes.
  Middleware ordering can otherwise create unverifiable behavior.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  Bound declared size before reading,
  read at most the configured limit plus one byte, hash the exact remaining
  stream bytes without decoding or normalization, close the original body, and
  restore an independent reader. Verification must be first; previously
  consumed bytes cannot be reconstructed and only the remaining bytes are
  authenticated.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestCaptureBodyPreservesExactBytesAndRestoresRequest`,
  `TestCaptureBodyPreservesEmptyCompressedTrailersAndPartialReads`,
  `TestCaptureBodyAfterPriorReadAuthenticatesOnlyRemainingBytes`, and body
  limit tests cover `CaptureBody`, `SignRequest`, and `VerifyRequest`.
  Reconsider if Go provides an earlier immutable raw-message boundary.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-009","title":"Exact request body, ordering, and ownership","status":"resolved","owner":"`webhook` maintainers","classification":"implementation-defined behavior","decision_scope":"defensive","specification":"RFC 9110 HTTP Semantics","version":"RFC 9110","source_authority":"rfc9110-source","section":"RFC 9110 Sections 6 and 8","requirement_strength":"not specified","issue":"RFC 9110 defines message content and Go `net/http` exposes a stream. Neither recovers bytes consumed by earlier middleware nor decides whether compressed content is decoded before authentication.","interpretations":["Hash decoded JSON, hash decompressed content, trust `Content-Length`, buffer without a cap, or hash exact bytes. Middleware ordering can otherwise create unverifiable behavior."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"Bound declared size before reading, read at most the configured limit plus one byte, hash the exact remaining stream bytes without decoding or normalization, close the original body, and restore an independent reader. Verification must be first; previously consumed bytes cannot be reconstructed and only the remaining bytes are authenticated.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestCaptureBodyPreservesExactBytesAndRestoresRequest","TestCaptureBodyBoundsUnknownLengthBeforeAllocation","TestCaptureBodyAfterPriorReadAuthenticatesOnlyRemainingBytes"],"fixture_evidence":[],"fuzz_evidence":[],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["CaptureBody","SignRequest","VerifyRequest"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc9110.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
implementation-defined behavior
defensive
RFC 9110 HTTP Semantics
RFC 9110
rfc9110-source
https://www.rfc-editor.org/rfc/rfc9110.txt
RFC 9110 Sections 6 and 8
not specified
RFC 9110 defines message content and Go `net/http` exposes a stream. Neither recovers bytes consumed by earlier middleware nor decides whether compressed content is decoded before authentication.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
Bound declared size before reading, read at most the configured limit plus one byte, hash the exact remaining stream bytes without decoding or normalization, close the original body, and restore an independent reader. Verification must be first; previously consumed bytes cannot be reconstructed and only the remaining bytes are authenticated.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Hash decoded JSON, hash decompressed content, trust `Content-Length`, buffer without a cap, or hash exact bytes. Middleware ordering can otherwise create unverifiable behavior.
TestCaptureBodyPreservesExactBytesAndRestoresRequest
TestCaptureBodyBoundsUnknownLengthBeforeAllocation
TestCaptureBodyAfterPriorReadAuthenticatesOnlyRemainingBytes
CaptureBody
SignRequest
VerifyRequest
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-010: Candidate verification and external errors

**Authoritative reference:** [RFC 2104](https://www.rfc-editor.org/rfc/rfc2104.html).

- **Status, owner, and classification:** `resolved`; maintainers; defensive
  verification and diagnostic policy.
- **Source and issue:** RFC 2104 defines MAC verification but does not define
  multi-key candidate processing, timing exposure across public key IDs, safe
  HTTP errors, or internal diagnostics.
- **Interpretations and peer behavior:** Fail on the first invalid candidate,
  expose exact causes, try configured candidates, or return a boolean. Detailed
  errors can disclose key lifecycle and parsing state.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  After strict set parsing, the
  verifier examines candidates in wire order and accepts the first active
  configured key whose MAC passes `hmac.Equal`. Invalid candidates do not alter
  later candidate state. External errors expose stable categories and the
  fixed message `webhook verification failed`; internal diagnostics remain
  separate and secret-safe. This is not a constant-runtime claim across public
  candidate metadata.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestVerifierSkipsMalformedAndInactiveSignaturesDeterministically`,
  `TestVerificationErrorMethodsAreNilSafe`,
  `TestMiddlewareReturnsOnlySafeFailureAndSkipsHandler`, and mutation tests
  cover `Verifier`, `VerificationError`, and middleware. Reconsider only if a
  stronger blinded key-selection design is required and measured.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-010","title":"Candidate verification and external errors","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"defensive","specification":"RFC 2104 HMAC and RFC 4231 HMAC-SHA-256/HMAC-SHA-512 vectors","version":"RFC 2104","source_authority":"rfc2104-source","section":"RFC 2104 Section 3","requirement_strength":"not specified","issue":"RFC 2104 defines MAC verification but does not define multi-key candidate processing, timing exposure across public key IDs, safe HTTP errors, or internal diagnostics.","interpretations":["Fail on the first invalid candidate, expose exact causes, try configured candidates, or return a boolean. Detailed errors can disclose key lifecycle and parsing state."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"After strict set parsing, the verifier examines candidates in wire order and accepts the first active configured key whose MAC passes `hmac.Equal`. Invalid candidates do not alter later candidate state. External errors expose stable categories and the fixed message `webhook verification failed`; internal diagnostics remain separate and secret-safe. This is not a constant-runtime claim across public candidate metadata.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestVerifierSkipsMalformedAndInactiveSignaturesDeterministically","TestVerificationErrorMethodsAreNilSafe"],"fixture_evidence":[],"fuzz_evidence":[],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["Verifier","VerificationError"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc2104.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
defensive
RFC 2104 HMAC and RFC 4231 HMAC-SHA-256/HMAC-SHA-512 vectors
RFC 2104
rfc2104-source
https://www.rfc-editor.org/rfc/rfc2104.txt
RFC 2104 Section 3
not specified
RFC 2104 defines MAC verification but does not define multi-key candidate processing, timing exposure across public key IDs, safe HTTP errors, or internal diagnostics.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
After strict set parsing, the verifier examines candidates in wire order and accepts the first active configured key whose MAC passes `hmac.Equal`. Invalid candidates do not alter later candidate state. External errors expose stable categories and the fixed message `webhook verification failed`; internal diagnostics remain separate and secret-safe. This is not a constant-runtime claim across public candidate metadata.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Fail on the first invalid candidate, expose exact causes, try configured candidates, or return a boolean. Detailed errors can disclose key lifecycle and parsing state.
TestVerifierSkipsMalformedAndInactiveSignaturesDeterministically
TestVerificationErrorMethodsAreNilSafe
Verifier
VerificationError
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-011: Replay identity and atomic storage

**Authoritative reference:** [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110.html).

- **Status, owner, and classification:** `resolved`; maintainers; defensive
  replay policy, not an exactly-once guarantee.
- **Source and issue:** HTTP and HMAC standards do not define webhook replay
  identity, tenant scoping, atomic persistence, TTL, or rotation behavior.
- **Interpretations and peer behavior:** Deduplicate by signature, nonce,
  provider event ID, key ID, or payload hash; check then insert; fail open on
  store outage; or delegate entirely to the application.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  Optional replay protection hashes a
  length-prefixed domain string, required namespace, and authenticated event ID
  with SHA-256. It deliberately excludes key ID so overlap signatures share one
  replay identity. `ReplayStore.CheckAndRecord` must atomically create only an
  absent key with expiry. Duplicate and backend error both fail closed. This
  prevents concurrent acceptance but does not claim exactly-once processing.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestVerifyAndRecordAtomicallyRejectsReplay`,
  `TestVerifyAndRecordHashesNamespacedReplayKey`,
  `TestReplayIdentitySurvivesSecretRotation`, and store adapter tests cover
  `ReplayStore`, `VerifyAndRecord`, and replay configuration. Reconsider for a
  versioned replay-key migration or a durable application transaction seam.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-011","title":"Replay identity and atomic storage","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"defensive","specification":"RFC 9110 HTTP Semantics","version":"RFC 9110","source_authority":"rfc9110-source","section":"RFC 9110 Sections 9 and 15","requirement_strength":"not specified","issue":"HTTP and HMAC standards do not define webhook replay identity, tenant scoping, atomic persistence, TTL, or rotation behavior.","interpretations":["Deduplicate by signature, nonce, provider event ID, key ID, or payload hash; check then insert; fail open on store outage; or delegate entirely to the application."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"Optional replay protection hashes a length-prefixed domain string, required namespace, and authenticated event ID with SHA-256. It deliberately excludes key ID so overlap signatures share one replay identity. `ReplayStore.CheckAndRecord` must atomically create only an absent key with expiry. Duplicate and backend error both fail closed. This prevents concurrent acceptance but does not claim exactly-once processing.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestVerifyAndRecordAtomicallyRejectsReplay","TestReplayIdentitySurvivesSecretRotation","TestVerifyAndRecordFailsClosedForMissingIDAndStoreFailure"],"fixture_evidence":[],"fuzz_evidence":[],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["ReplayStore","Verifier.VerifyAndRecord"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc9110.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
defensive
RFC 9110 HTTP Semantics
RFC 9110
rfc9110-source
https://www.rfc-editor.org/rfc/rfc9110.txt
RFC 9110 Sections 9 and 15
not specified
HTTP and HMAC standards do not define webhook replay identity, tenant scoping, atomic persistence, TTL, or rotation behavior.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
Optional replay protection hashes a length-prefixed domain string, required namespace, and authenticated event ID with SHA-256. It deliberately excludes key ID so overlap signatures share one replay identity. `ReplayStore.CheckAndRecord` must atomically create only an absent key with expiry. Duplicate and backend error both fail closed. This prevents concurrent acceptance but does not claim exactly-once processing.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Deduplicate by signature, nonce, provider event ID, key ID, or payload hash; check then insert; fail open on store outage; or delegate entirely to the application.
TestVerifyAndRecordAtomicallyRejectsReplay
TestReplayIdentitySurvivesSecretRotation
TestVerifyAndRecordFailsClosedForMissingIDAndStoreFailure
ReplayStore
Verifier.VerifyAndRecord
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-012: Event ID extraction occurs after authentication

**Authoritative reference:** [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110.html).

- **Status, owner, and classification:** `resolved`; maintainers; defensive
  sequencing and application extension policy.
- **Source and issue:** No standard defines where a generic event ID resides.
  Extracting attacker-controlled JSON or headers before authentication can
  consume resources, leak parser behavior, or poison replay state.
- **Interpretations and peer behavior:** Require one fixed header, decode JSON
  before verification, let every handler deduplicate, or inject an extractor
  after authentication.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  `VerifyRequest` authenticates strict
  fields and exact body first, then invokes the configured extractor, then
  atomically records replay state. `HeaderEventID` requires one bounded UTF-8
  field. Extractor errors map to a safe missing-ID category and unauthenticated
  requests cannot touch replay storage.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestVerifyRequestExtractsEventIDAfterAuthentication`,
  `TestHeaderEventIDRejectsDuplicateAndOversizedValues`, and
  `TestVerifyRequestAndReplayErrorPaths` cover `EventIDExtractor`,
  `HeaderEventID`, and `VerifyRequest`. Reconsider only for an isolated provider
  profile with an authoritative event-ID location.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-012","title":"Event ID extraction occurs after authentication","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"defensive","specification":"RFC 9110 HTTP Semantics","version":"RFC 9110","source_authority":"rfc9110-source","section":"RFC 9110 Sections 5-8","requirement_strength":"not specified","issue":"No standard defines where a generic event ID resides. Extracting attacker-controlled JSON or headers before authentication can consume resources, leak parser behavior, or poison replay state.","interpretations":["Require one fixed header, decode JSON before verification, let every handler deduplicate, or inject an extractor after authentication."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"`VerifyRequest` authenticates strict fields and exact body first, then invokes the configured extractor, then atomically records replay state. `HeaderEventID` requires one bounded UTF-8 field. Extractor errors map to a safe missing-ID category and unauthenticated requests cannot touch replay storage.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestVerifyRequestExtractsEventIDAfterAuthentication","TestHeaderEventIDRejectsDuplicateAndOversizedValues"],"fixture_evidence":[],"fuzz_evidence":[],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["EventIDExtractor","VerifyRequest"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc9110.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
defensive
RFC 9110 HTTP Semantics
RFC 9110
rfc9110-source
https://www.rfc-editor.org/rfc/rfc9110.txt
RFC 9110 Sections 5-8
not specified
No standard defines where a generic event ID resides. Extracting attacker-controlled JSON or headers before authentication can consume resources, leak parser behavior, or poison replay state.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
`VerifyRequest` authenticates strict fields and exact body first, then invokes the configured extractor, then atomically records replay state. `HeaderEventID` requires one bounded UTF-8 field. Extractor errors map to a safe missing-ID category and unauthenticated requests cannot touch replay storage.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Require one fixed header, decode JSON before verification, let every handler deduplicate, or inject an extractor after authentication.
TestVerifyRequestExtractsEventIDAfterAuthentication
TestHeaderEventIDRejectsDuplicateAndOversizedValues
EventIDExtractor
VerifyRequest
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-013: Envelope is CloudEvents-shaped but not CloudEvents

**Authoritative reference:** [CloudEvents 1.0.2](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md).

- **Status, owner, and classification:** `resolved`; maintainers; local JSON
  wire policy with an explicit non-conformance boundary.
- **Source and issue:** CloudEvents 1.0.2 defines context attributes, extension
  attributes, data content, and constraints. `Envelope` uses familiar field
  names and emits `specversion: "1.0"` but also requires time, permits an
  unvalidated source string, and nests arbitrary metadata rather than exposing
  CloudEvents extension attributes.
- **Interpretations and peer behavior:** Claim structured CloudEvents JSON,
  remove familiar names, validate the complete CloudEvents model, or preserve
  the existing small local envelope while stating the boundary.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  The output is a deterministic local
  v1 JSON envelope, not a CloudEvents implementation or interoperability claim.
  It requires ID, type, source, nonzero time, and valid JSON data; emits UTC
  RFC3339Nano time and `application/json`; preserves raw data; and orders fields
  through a fixed Go struct. Metadata is a nested local object.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestEnvelopeMarshalIsDeterministicAndPreservesData`,
  `TestEnvelopeRejectsInvalidRequiredFieldsAndData`, and `FuzzEnvelope` cover
  `Envelope.MarshalJSON`. Reconsider by moving true CloudEvents support to the
  existing `cloudevents` package or by introducing a distinctly versioned
  envelope without misleading overlap.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-013","title":"Envelope is CloudEvents-shaped but not CloudEvents","status":"resolved","owner":"`webhook` maintainers","classification":"interoperability policy","decision_scope":"application-policy","specification":"CloudEvents 1.0.2 comparison boundary","version":"CloudEvents 1.0.2","source_authority":"cloudevents-core-v1.0.2","section":"CloudEvents 1.0.2 Sections 1-3","requirement_strength":"not specified","issue":"CloudEvents 1.0.2 defines context attributes, extension attributes, data content, and constraints. `Envelope` uses familiar field names and emits `specversion: \"1.0\"` but also requires time, permits an unvalidated source string, and nests arbitrary metadata rather than exposing CloudEvents extension attributes.","interpretations":["Claim structured CloudEvents JSON, remove familiar names, validate the complete CloudEvents model, or preserve the existing small local envelope while stating the boundary."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"The output is a deterministic local v1 JSON envelope, not a CloudEvents implementation or interoperability claim. It requires ID, type, source, nonzero time, and valid JSON data; emits UTC RFC3339Nano time and `application/json`; preserves raw data; and orders fields through a fixed Go struct. Metadata is a nested local object.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestEnvelopeMarshalIsDeterministicAndPreservesData","TestEnvelopeRejectsInvalidRequiredFieldsAndData"],"fixture_evidence":[],"fuzz_evidence":["FuzzEnvelope"],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["Envelope"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://raw.githubusercontent.com/cloudevents/spec/fc1f6f31f5f011a72183f1bcea20c987cb683ade/cloudevents/spec.md

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
interoperability policy
application-policy
CloudEvents 1.0.2 comparison boundary
CloudEvents 1.0.2
cloudevents-core-v1.0.2
https://raw.githubusercontent.com/cloudevents/spec/fc1f6f31f5f011a72183f1bcea20c987cb683ade/cloudevents/spec.md
CloudEvents 1.0.2 Sections 1-3
not specified
CloudEvents 1.0.2 defines context attributes, extension attributes, data content, and constraints. `Envelope` uses familiar field names and emits `specversion: "1.0"` but also requires time, permits an unvalidated source string, and nests arbitrary metadata rather than exposing CloudEvents extension attributes.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
The output is a deterministic local v1 JSON envelope, not a CloudEvents implementation or interoperability claim. It requires ID, type, source, nonzero time, and valid JSON data; emits UTC RFC3339Nano time and `application/json`; preserves raw data; and orders fields through a fixed Go struct. Metadata is a nested local object.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Claim structured CloudEvents JSON, remove familiar names, validate the complete CloudEvents model, or preserve the existing small local envelope while stating the boundary.
TestEnvelopeMarshalIsDeterministicAndPreservesData
TestEnvelopeRejectsInvalidRequiredFieldsAndData
FuzzEnvelope
Envelope
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-014: Outbound method, content type, and idempotency field

**Authoritative reference:** [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110.html).

- **Status, owner, and classification:** `resolved`; maintainers; application
  delivery profile on top of RFC 9110.
- **Source and issue:** HTTP permits many methods and media types and does not
  standardize receiver idempotency semantics. A generic sender still needs one
  deterministic request profile.
- **Interpretations and peer behavior:** Preserve arbitrary caller method and
  content type, infer JSON, use PUT for idempotency, or define a fixed POST
  profile. Webhook receivers overwhelmingly vary by provider.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  `Deliverer` always emits POST with
  `Content-Type: application/json`, preserves other cloned caller headers, and
  emits `Idempotency-Key` only when explicitly supplied. Both fixed field
  values are signed. The package does not claim that a receiver honors that
  application-defined idempotency field.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestDeliverRetriesNetworkFailureAndStopsAtBound`,
  `TestDeliverPreservesRequestStatusAndResponseBoundaries`, and fixed-header
  mutation tests cover `DeliveryRequest` and `Deliverer`. Reconsider with an
  explicit configurable request profile and corresponding signature version.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-014","title":"Outbound method, content type, and idempotency field","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"application-policy","specification":"RFC 9110 HTTP Semantics","version":"RFC 9110","source_authority":"rfc9110-source","section":"RFC 9110 Sections 7-9","requirement_strength":"not specified","issue":"HTTP permits many methods and media types and does not standardize receiver idempotency semantics. A generic sender still needs one deterministic request profile.","interpretations":["Preserve arbitrary caller method and content type, infer JSON, use PUT for idempotency, or define a fixed POST profile. Webhook receivers overwhelmingly vary by provider."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"`Deliverer` always emits POST with `Content-Type: application/json`, preserves other cloned caller headers, and emits `Idempotency-Key` only when explicitly supplied. Both fixed field values are signed. The package does not claim that a receiver honors that application-defined idempotency field.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestDeliverPreservesRequestStatusAndResponseBoundaries","TestVerifyRequestRejectsMutationOfFixedSignedHeaders"],"fixture_evidence":[],"fuzz_evidence":[],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["DeliveryRequest","Deliverer"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc9110.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
application-policy
RFC 9110 HTTP Semantics
RFC 9110
rfc9110-source
https://www.rfc-editor.org/rfc/rfc9110.txt
RFC 9110 Sections 7-9
not specified
HTTP permits many methods and media types and does not standardize receiver idempotency semantics. A generic sender still needs one deterministic request profile.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
`Deliverer` always emits POST with `Content-Type: application/json`, preserves other cloned caller headers, and emits `Idempotency-Key` only when explicitly supplied. Both fixed field values are signed. The package does not claim that a receiver honors that application-defined idempotency field.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Preserve arbitrary caller method and content type, infer JSON, use PUT for idempotency, or define a fixed POST profile. Webhook receivers overwhelmingly vary by provider.
TestDeliverPreservesRequestStatusAndResponseBoundaries
TestVerifyRequestRejectsMutationOfFixedSignedHeaders
DeliveryRequest
Deliverer
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-015: Retryable outcomes and Retry-After

**Authoritative reference:** [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110.html).

- **Status, owner, and classification:** `resolved`; maintainers; RFC 9110
  interpretation plus bounded application policy.
- **Source and issue:** RFC 9110 defines `Retry-After` and status semantics but
  does not require clients to retry, define transport failure treatment, or
  select a webhook status allowlist.
- **Interpretations and peer behavior:** Retry all 4xx/5xx, only 429/503, every
  transport error, or caller-configured statuses. Clients also disagree on
  invalid, past, overflowing, and excessive `Retry-After` values.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  Retry transport failures and exactly
  408, 425, 429, 500, 502, 503, and 504. Any 2xx succeeds; every other status is
  terminal. Parse nonnegative delta-seconds or a future HTTP-date, cap at
  `MaxDelay`, and otherwise use capped exponential delay. Overflow saturates
  before capping; cancellation interrupts backoff.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestDeliverRetriesRetryableStatusAndHonorsRetryAfter`,
  `TestRetryAfterSupportsHTTPDateAndCapsDelay`,
  `TestRetryPolicyCapsOverflowingRetryAfterSeconds`, and cancellation tests
  cover `RetryPolicy` and `Deliver`. Reconsider when receiver contracts require
  an explicit configurable classifier rather than changing this list silently.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-015","title":"Retryable outcomes and Retry-After","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"application-policy","specification":"RFC 9110 HTTP Semantics","version":"RFC 9110","source_authority":"rfc9110-source","section":"RFC 9110 Sections 10.2.3 and 15","requirement_strength":"not specified","issue":"RFC 9110 defines `Retry-After` and status semantics but does not require clients to retry, define transport failure treatment, or select a webhook status allowlist.","interpretations":["Retry all 4xx/5xx, only 429/503, every transport error, or caller-configured statuses. Clients also disagree on invalid, past, overflowing, and excessive `Retry-After` values."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"Retry transport failures and exactly 408, 425, 429, 500, 502, 503, and 504. Any 2xx succeeds; every other status is terminal. Parse nonnegative delta-seconds or a future HTTP-date, cap at `MaxDelay`, and otherwise use capped exponential delay. Overflow saturates before capping; cancellation interrupts backoff.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestDeliverRetriesRetryableStatusAndHonorsRetryAfter","TestRetryAfterSupportsHTTPDateAndCapsDelay","TestRetryPolicyCapsOverflowingRetryAfterSeconds"],"fixture_evidence":[],"fuzz_evidence":[],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["RetryPolicy","DeliveryAttempt"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc9110.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
application-policy
RFC 9110 HTTP Semantics
RFC 9110
rfc9110-source
https://www.rfc-editor.org/rfc/rfc9110.txt
RFC 9110 Sections 10.2.3 and 15
not specified
RFC 9110 defines `Retry-After` and status semantics but does not require clients to retry, define transport failure treatment, or select a webhook status allowlist.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
Retry transport failures and exactly 408, 425, 429, 500, 502, 503, and 504. Any 2xx succeeds; every other status is terminal. Parse nonnegative delta-seconds or a future HTTP-date, cap at `MaxDelay`, and otherwise use capped exponential delay. Overflow saturates before capping; cancellation interrupts backoff.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Retry all 4xx/5xx, only 429/503, every transport error, or caller-configured statuses. Clients also disagree on invalid, past, overflowing, and excessive `Retry-After` values.
TestDeliverRetriesRetryableStatusAndHonorsRetryAfter
TestRetryAfterSupportsHTTPDateAndCapsDelay
TestRetryPolicyCapsOverflowingRetryAfterSeconds
RetryPolicy
DeliveryAttempt
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-016: Retry ownership and ambiguous receipt

**Authoritative reference:** [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110.html).

- **Status, owner, and classification:** `resolved`; maintainers; delivery
  lifecycle and idempotency safety policy.
- **Source and issue:** HTTP cannot reveal whether a timed-out request caused a
  side effect. Layered HTTP, queue, and outbox retries can multiply attempts.
- **Interpretations and peer behavior:** Retry regardless, require an
  idempotency key, let every layer retry, or make durable infrastructure the
  sole retry owner.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  `Deliver` forces one attempt when no
  explicit idempotency key exists. `DeliverOnce` always performs one attempt
  and is mandatory for queue/outbox consumers, leaving durable retry ownership
  outside the core. This reduces duplicate risk but still makes no exactly-once
  claim about a remote endpoint.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestDeliverOnceDisablesInternalRetries`,
  `TestHandleUsesSingleDeliveryAttempt`,
  `TestPublisherPerformsSingleAttemptForRelay`, and ambiguous failure tests
  cover `Deliver`, `DeliverOnce`, and adapters. Reconsider only with a durable
  cross-boundary protocol that proves receiver deduplication.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-016","title":"Retry ownership and ambiguous receipt","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"defensive","specification":"RFC 9110 HTTP Semantics","version":"RFC 9110","source_authority":"rfc9110-source","section":"RFC 9110 Sections 9 and 15","requirement_strength":"not specified","issue":"HTTP cannot reveal whether a timed-out request caused a side effect. Layered HTTP, queue, and outbox retries can multiply attempts.","interpretations":["Retry regardless, require an idempotency key, let every layer retry, or make durable infrastructure the sole retry owner."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"`Deliver` forces one attempt when no explicit idempotency key exists. `DeliverOnce` always performs one attempt and is mandatory for queue/outbox consumers, leaving durable retry ownership outside the core. This reduces duplicate risk but still makes no exactly-once claim about a remote endpoint.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestDeliverRetriesNetworkFailureAndStopsAtBound","TestDeliverOnceDisablesInternalRetries","TestHandleUsesSingleDeliveryAttempt","TestPublisherPerformsSingleAttemptForRelay"],"fixture_evidence":[],"fuzz_evidence":[],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["Deliverer","DeliveryRequest"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc9110.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
defensive
RFC 9110 HTTP Semantics
RFC 9110
rfc9110-source
https://www.rfc-editor.org/rfc/rfc9110.txt
RFC 9110 Sections 9 and 15
not specified
HTTP cannot reveal whether a timed-out request caused a side effect. Layered HTTP, queue, and outbox retries can multiply attempts.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
`Deliver` forces one attempt when no explicit idempotency key exists. `DeliverOnce` always performs one attempt and is mandatory for queue/outbox consumers, leaving durable retry ownership outside the core. This reduces duplicate risk but still makes no exactly-once claim about a remote endpoint.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Retry regardless, require an idempotency key, let every layer retry, or make durable infrastructure the sole retry owner.
TestDeliverRetriesNetworkFailureAndStopsAtBound
TestDeliverOnceDisablesInternalRetries
TestHandleUsesSingleDeliveryAttempt
TestPublisherPerformsSingleAttemptForRelay
Deliverer
DeliveryRequest
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-017: Response bounds and failure classification

**Authoritative reference:** [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110.html).

- **Status, owner, and classification:** `resolved`; maintainers; defensive
  resource and result policy.
- **Source and issue:** RFC 9110 does not define how much response content a
  webhook sender retains, whether body close failures matter, or how attempts
  map to stable application classifications.
- **Interpretations and peer behavior:** Ignore response bodies, read without a
  bound, truncate silently, expose transport errors verbatim, or preserve a
  bounded result and stable category.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  Require positive request and response
  bounds that leave room for a sentinel byte. Read at most response limit plus
  one, close on every path, reject missing bodies and read/close failures, and
  classify each attempt as none, retryable, terminal, or exhausted. Diagnostic
  strings remain fixed and payloads or sensitive fields are not recorded.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestDeliverBoundsResponseBody`, `TestReadResponseRejectsMissingReadAndCloseFailures`,
  `TestDeliverClosesResponseReturnedWithTransportError`, and classification
  tests cover `DeliveryAttempt`, `DeliveryResult`, and `Deliver`. Reconsider if
  streaming response consumption becomes an explicit separate API.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-017","title":"Response bounds and failure classification","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"defensive","specification":"RFC 9110 HTTP Semantics","version":"RFC 9110","source_authority":"rfc9110-source","section":"RFC 9110 Sections 6 and 15","requirement_strength":"not specified","issue":"RFC 9110 does not define how much response content a webhook sender retains, whether body close failures matter, or how attempts map to stable application classifications.","interpretations":["Ignore response bodies, read without a bound, truncate silently, expose transport errors verbatim, or preserve a bounded result and stable category."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"Require positive request and response bounds that leave room for a sentinel byte. Read at most response limit plus one, close on every path, reject missing bodies and read/close failures, and classify each attempt as none, retryable, terminal, or exhausted. Diagnostic strings remain fixed and payloads or sensitive fields are not recorded.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestDeliverBoundsResponseBody","TestReadResponseRejectsMissingReadAndCloseFailures","TestDeliverClassifiesCanceledTransportAndExhaustedStatus"],"fixture_evidence":[],"fuzz_evidence":[],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["DeliveryResult","DeliveryAttempt"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc9110.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
defensive
RFC 9110 HTTP Semantics
RFC 9110
rfc9110-source
https://www.rfc-editor.org/rfc/rfc9110.txt
RFC 9110 Sections 6 and 15
not specified
RFC 9110 does not define how much response content a webhook sender retains, whether body close failures matter, or how attempts map to stable application classifications.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
Require positive request and response bounds that leave room for a sentinel byte. Read at most response limit plus one, close on every path, reject missing bodies and read/close failures, and classify each attempt as none, retryable, terminal, or exhausted. Diagnostic strings remain fixed and payloads or sensitive fields are not recorded.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Ignore response bodies, read without a bound, truncate silently, expose transport errors verbatim, or preserve a bounded result and stable category.
TestDeliverBoundsResponseBody
TestReadResponseRejectsMissingReadAndCloseFailures
TestDeliverClassifiesCanceledTransportAndExhaustedStatus
DeliveryResult
DeliveryAttempt
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-018: Dead letters and operator replay are hooks

**Authoritative reference:** [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110.html).

- **Status, owner, and classification:** `resolved`; maintainers; application
  lifecycle boundary.
- **Source and issue:** HTTP defines neither dead-letter persistence nor an
  operator replay audit. Implementing storage or a queue in the core would
  duplicate `outbox` and `queue` ownership.
- **Interpretations and peer behavior:** Persist internally, drop terminal
  results, expose callbacks, or require one external orchestration package.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  Terminal and exhausted deliveries
  invoke an optional bounded-result `DeadLetterFunc`. Operator `Replay` invokes
  its audit hook before creating a new delivery ID and never reuses an attempt
  ID. Hook failures are returned; observer panics are contained, while
  lifecycle hook ownership remains explicit. The core stores nothing.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestDeliverDoesNotRetryTerminalStatusAndDeadLetters`,
  `TestDeliveryPreservesDeadLetterHookFailure`, and
  `TestReplayAuditsBeforeStartingNewDelivery` cover `DeadLetterFunc`,
  `ReplayHook`, and `Replay`. Reconsider only if a durable adapter contract
  cannot express a required atomic transition.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-018","title":"Dead letters and operator replay are hooks","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"application-policy","specification":"RFC 9110 HTTP Semantics","version":"RFC 9110","source_authority":"rfc9110-source","section":"RFC 9110 Sections 9 and 15","requirement_strength":"not specified","issue":"HTTP defines neither dead-letter persistence nor an operator replay audit. Implementing storage or a queue in the core would duplicate `outbox` and `queue` ownership.","interpretations":["Persist internally, drop terminal results, expose callbacks, or require one external orchestration package."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"Terminal and exhausted deliveries invoke an optional bounded-result `DeadLetterFunc`. Operator `Replay` invokes its audit hook before creating a new delivery ID and never reuses an attempt ID. Hook failures are returned; observer panics are contained, while lifecycle hook ownership remains explicit. The core stores nothing.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestDeliverDoesNotRetryTerminalStatusAndDeadLetters","TestDeliveryPreservesDeadLetterHookFailure","TestReplayAuditsBeforeStartingNewDelivery"],"fixture_evidence":[],"fuzz_evidence":[],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["DeadLetterHook","ReplayHook","Deliverer"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc9110.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
application-policy
RFC 9110 HTTP Semantics
RFC 9110
rfc9110-source
https://www.rfc-editor.org/rfc/rfc9110.txt
RFC 9110 Sections 9 and 15
not specified
HTTP defines neither dead-letter persistence nor an operator replay audit. Implementing storage or a queue in the core would duplicate `outbox` and `queue` ownership.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
Terminal and exhausted deliveries invoke an optional bounded-result `DeadLetterFunc`. Operator `Replay` invokes its audit hook before creating a new delivery ID and never reuses an attempt ID. Hook failures are returned; observer panics are contained, while lifecycle hook ownership remains explicit. The core stores nothing.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Persist internally, drop terminal results, expose callbacks, or require one external orchestration package.
TestDeliverDoesNotRetryTerminalStatusAndDeadLetters
TestDeliveryPreservesDeadLetterHookFailure
TestReplayAuditsBeforeStartingNewDelivery
DeadLetterHook
ReplayHook
Deliverer
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-019: Endpoint URL syntax and scheme policy

**Authoritative reference:** [RFC 3986](https://www.rfc-editor.org/rfc/rfc3986.html).

- **Status, owner, and classification:** `resolved`; maintainers; RFC 3986
  parsing with defensive SSRF policy.
- **Source and issue:** RFC 3986 permits URI forms and components that are not
  safe outbound webhook endpoints. Go `net/url` accepts opaque URLs, userinfo,
  fragments, non-ASCII hosts, trailing dots, and varied port spellings.
- **Interpretations and peer behavior:** Accept anything `url.Parse` accepts,
  allow HTTP by default, normalize risky forms, or reject unless explicitly
  enabled. SSRF filters often disagree on canonical host treatment.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  Require an absolute hierarchical URL
  with host; exact lowercase `https` is default and lowercase `http` is opt-in.
  Reject opaque forms, userinfo, fragments, empty hosts, non-ASCII hosts,
  trailing-dot hosts, malformed ports, and ports outside 1-65535. Do not
  silently IDNA-convert or rewrite attacker input.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestSSRFPolicyRejectsUnsafeURLsAndAddresses`,
  `TestSSRFPolicyConfigurationAndAddressFailures`, and `FuzzSSRFPolicy` cover
  `SSRFPolicyConfig`, `NewSSRFPolicy`, and `Validate`. Reconsider if explicit
  IDNA policy is added with hostile normalization vectors.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-019","title":"Endpoint URL syntax and scheme policy","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"defensive","specification":"RFC 3986 URI Generic Syntax","version":"RFC 3986","source_authority":"rfc3986-source","section":"RFC 3986 Sections 3-5","requirement_strength":"not specified","issue":"RFC 3986 permits URI forms and components that are not safe outbound webhook endpoints. Go `net/url` accepts opaque URLs, userinfo, fragments, non-ASCII hosts, trailing dots, and varied port spellings.","interpretations":["Accept anything `url.Parse` accepts, allow HTTP by default, normalize risky forms, or reject unless explicitly enabled. SSRF filters often disagree on canonical host treatment."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"Require an absolute hierarchical URL with host; exact lowercase `https` is default and lowercase `http` is opt-in. Reject opaque forms, userinfo, fragments, empty hosts, non-ASCII hosts, trailing-dot hosts, malformed ports, and ports outside 1-65535. Do not silently IDNA-convert or rewrite attacker input.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestSSRFPolicyRejectsUnsafeURLsAndAddresses","TestSSRFPolicyAllowsExplicitPrefixOnlyWhenConfigured"],"fixture_evidence":[],"fuzz_evidence":["FuzzSSRFPolicy"],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["SSRFPolicy"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc3986.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
defensive
RFC 3986 URI Generic Syntax
RFC 3986
rfc3986-source
https://www.rfc-editor.org/rfc/rfc3986.txt
RFC 3986 Sections 3-5
not specified
RFC 3986 permits URI forms and components that are not safe outbound webhook endpoints. Go `net/url` accepts opaque URLs, userinfo, fragments, non-ASCII hosts, trailing dots, and varied port spellings.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
Require an absolute hierarchical URL with host; exact lowercase `https` is default and lowercase `http` is opt-in. Reject opaque forms, userinfo, fragments, empty hosts, non-ASCII hosts, trailing-dot hosts, malformed ports, and ports outside 1-65535. Do not silently IDNA-convert or rewrite attacker input.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Accept anything `url.Parse` accepts, allow HTTP by default, normalize risky forms, or reject unless explicitly enabled. SSRF filters often disagree on canonical host treatment.
TestSSRFPolicyRejectsUnsafeURLsAndAddresses
TestSSRFPolicyAllowsExplicitPrefixOnlyWhenConfigured
FuzzSSRFPolicy
SSRFPolicy
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-020: Special-purpose address registry policy

**Authoritative reference:** [IANA IPv4 special-purpose registry](https://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xhtml).

- **Status, owner, and classification:** `resolved`; maintainers; defensive
  network policy based on pinned IANA registry snapshots.
- **Source and issue:** IANA IPv4 and IPv6 special-purpose registries classify
  ranges with several properties, while Go address predicates do not reject
  every documentation, benchmark, shared, reserved, or future-use range needed
  by a conservative internet-delivery policy.
- **Interpretations and peer behavior:** Allow every global-unicast address,
  block only RFC 1918, mirror registry forwarding flags, or deny a conservative
  reviewed set. Library filters commonly miss mapped IPv4 and mixed DNS answers.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  Unmap IPv4-mapped IPv6; deny invalid,
  private, loopback, link-local, multicast, unspecified, and non-global-unicast
  addresses; additionally deny the explicit pinned special-purpose prefixes in
  `reservedPrefixes`. Caller deny prefixes win over caller allow prefixes.
  Explicit allow prefixes are narrow operator exceptions. Every DNS answer
  must pass and answer count is bounded.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestSSRFPolicyRejectsUnsafeURLsAndAddresses`,
  `TestSSRFPolicyRejectsMixedAndOversizedDNSAnswers`, and
  `TestSSRFPolicyAllowsExplicitPrefixOnlyWhenConfigured` cover address policy.
  Reconsider whenever either pinned IANA registry digest changes; review the
  range diff before changing runtime acceptance.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-020","title":"Special-purpose address registry policy","status":"resolved","owner":"`webhook` maintainers","classification":"interoperability policy","decision_scope":"defensive","specification":"IANA IPv4 and IPv6 Special-Purpose Address Registries at 2025-10-09","version":"IANA special-purpose registries 2025-10-09","source_authority":"iana-ipv4-source","section":"IANA IPv4 and IPv6 special-purpose registries","requirement_strength":"not specified","issue":"IANA IPv4 and IPv6 special-purpose registries classify ranges with several properties, while Go address predicates do not reject every documentation, benchmark, shared, reserved, or future-use range needed by a conservative internet-delivery policy.","interpretations":["Allow every global-unicast address, block only RFC 1918, mirror registry forwarding flags, or deny a conservative reviewed set. Library filters commonly miss mapped IPv4 and mixed DNS answers."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"Unmap IPv4-mapped IPv6; deny invalid, private, loopback, link-local, multicast, unspecified, and non-global-unicast addresses; additionally deny the explicit pinned special-purpose prefixes in `reservedPrefixes`. Caller deny prefixes win over caller allow prefixes. Explicit allow prefixes are narrow operator exceptions. Every DNS answer must pass and answer count is bounded.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestSSRFPolicyRejectsUnsafeURLsAndAddresses","TestSSRFPolicyRejectsMixedAndOversizedDNSAnswers"],"fixture_evidence":[],"fuzz_evidence":["FuzzSSRFPolicy"],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["SSRFPolicy"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xml

Additional authoritative source: `{"id":"iana-ipv6-source","version":"IANA special-purpose registries 2025-10-09","url":"https://www.iana.org/assignments/iana-ipv6-special-registry/iana-ipv6-special-registry.xml","specifications":["IANA IPv4 and IPv6 Special-Purpose Address Registries at 2025-10-09"]}`

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
interoperability policy
defensive
IANA IPv4 and IPv6 Special-Purpose Address Registries at 2025-10-09
IANA special-purpose registries 2025-10-09
iana-ipv4-source
https://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xml
IANA IPv4 and IPv6 special-purpose registries
not specified
IANA IPv4 and IPv6 special-purpose registries classify ranges with several properties, while Go address predicates do not reject every documentation, benchmark, shared, reserved, or future-use range needed by a conservative internet-delivery policy.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
Unmap IPv4-mapped IPv6; deny invalid, private, loopback, link-local, multicast, unspecified, and non-global-unicast addresses; additionally deny the explicit pinned special-purpose prefixes in `reservedPrefixes`. Caller deny prefixes win over caller allow prefixes. Explicit allow prefixes are narrow operator exceptions. Every DNS answer must pass and answer count is bounded.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Allow every global-unicast address, block only RFC 1918, mirror registry forwarding flags, or deny a conservative reviewed set. Library filters commonly miss mapped IPv4 and mixed DNS answers.
TestSSRFPolicyRejectsUnsafeURLsAndAddresses
TestSSRFPolicyRejectsMixedAndOversizedDNSAnswers
FuzzSSRFPolicy
SSRFPolicy
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-021: DNS rebinding, redirects, proxies, and transport

**Authoritative reference:** [Go 1.26.6 net/http source](https://cs.opensource.google/go/go/+/refs/tags/go1.26.6:src/net/http/).

- **Status, owner, and classification:** `resolved`; maintainers; defensive Go
  HTTP transport policy.
- **Source and issue:** Go `net/http` owns redirect, proxy, DNS, connection
  reuse, and protocol behavior. URL validation before a request alone does not
  prevent DNS rebinding, redirect pivots, or environment proxy bypass.
- **Interpretations and peer behavior:** Validate only the original URL, follow
  redirects with revalidation, trust environment proxies, pin one DNS answer,
  or own a direct dial path. Generic clients often validate too early.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  Validate immediately before every
  attempt and again in the transport; re-resolve at dial time; validate every
  answer; dial only validated addresses; disable environment proxies; return
  redirects without following; and disable automatic HTTP/2 on this custom
  transport so its direct validated dial ownership remains explicit. TLS
  certificate and hostname verification remain with `net/http`.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestSecureHTTPClientRevalidatesDNSAtDialTime`,
  `TestSecureHTTPClientRejectsRedirectWithoutContactingTarget`,
  `TestSecureHTTPClientConfigurationAndTransportFailures`, and transport
  configuration tests cover `NewSecureHTTPClient`. Reconsider if Go exposes a
  protocol-independent validated-address dial contract with equivalent proof.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-021","title":"DNS rebinding, redirects, proxies, and transport","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"defensive","specification":"Go 1.26.6 cryptography, HTTP, URL, address, time, and encoding contracts","version":"Go 1.26.6","source_authority":"go-net-http-source","section":"net/http transport, redirect, proxy, and dialing contracts","requirement_strength":"not specified","issue":"Go `net/http` owns redirect, proxy, DNS, connection reuse, and protocol behavior. URL validation before a request alone does not prevent DNS rebinding, redirect pivots, or environment proxy bypass.","interpretations":["Validate only the original URL, follow redirects with revalidation, trust environment proxies, pin one DNS answer, or own a direct dial path. Generic clients often validate too early."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"Validate immediately before every attempt and again in the transport; re-resolve at dial time; validate every answer; dial only validated addresses; disable environment proxies; return redirects without following; and disable automatic HTTP/2 on this custom transport so its direct validated dial ownership remains explicit. TLS certificate and hostname verification remain with `net/http`.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestSecureHTTPClientRevalidatesDNSAtDialTime","TestSecureHTTPClientRejectsRedirectWithoutContactingTarget","TestSecureHTTPClientConfigurationAndTransportFailures"],"fixture_evidence":[],"fuzz_evidence":["FuzzSSRFPolicy"],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["NewSecureHTTPClient","SSRFPolicy"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://raw.githubusercontent.com/golang/go/go1.26.6/src/net/http/server.go

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
defensive
Go 1.26.6 cryptography, HTTP, URL, address, time, and encoding contracts
Go 1.26.6
go-net-http-source
https://raw.githubusercontent.com/golang/go/go1.26.6/src/net/http/server.go
net/http transport, redirect, proxy, and dialing contracts
not specified
Go `net/http` owns redirect, proxy, DNS, connection reuse, and protocol behavior. URL validation before a request alone does not prevent DNS rebinding, redirect pivots, or environment proxy bypass.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
Validate immediately before every attempt and again in the transport; re-resolve at dial time; validate every answer; dial only validated addresses; disable environment proxies; return redirects without following; and disable automatic HTTP/2 on this custom transport so its direct validated dial ownership remains explicit. TLS certificate and hostname verification remain with `net/http`.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Validate only the original URL, follow redirects with revalidation, trust environment proxies, pin one DNS answer, or own a direct dial path. Generic clients often validate too early.
TestSecureHTTPClientRevalidatesDNSAtDialTime
TestSecureHTTPClientRejectsRedirectWithoutContactingTarget
TestSecureHTTPClientConfigurationAndTransportFailures
FuzzSSRFPolicy
NewSecureHTTPClient
SSRFPolicy
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-022: Fan-out concurrency and ordering

**Authoritative reference:** [Go memory model](https://go.dev/ref/mem).

- **Status, owner, and classification:** `resolved`; maintainers; bounded
  orchestration policy.
- **Source and issue:** Neither HTTP nor webhook conventions define fan-out
  concurrency, result ordering, cancellation, or durability. One goroutine per
  endpoint is an unbounded resource risk.
- **Interpretations and peer behavior:** Spawn freely, run serially, return
  completion order, stop on first failure, or use a fixed worker set.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  `FanOut` rejects input beyond
  `MaxFanOut`, requires a positive worker limit, uses a fixed worker bound, and
  returns one result per input in input order. Cancellation prevents new useful
  work but does not claim to interrupt a remote side effect already in flight.
  Fan-out is not a durable queue.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestFanOutBoundsConcurrencyAndPreservesResultOrder`,
  `TestFanOutRejectsUnboundedInputs`, and identifier boundary tests cover
  `FanOut`. Reconsider only with an explicit streaming-result API or durable
  orchestration adapter.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-022","title":"Fan-out concurrency and ordering","status":"resolved","owner":"`webhook` maintainers","classification":"implementation-defined behavior","decision_scope":"application-policy","specification":"Go 1.26.6 cryptography, HTTP, URL, address, time, and encoding contracts","version":"Go 1.26.6","source_authority":"go-language-source","section":"Go language concurrency and memory-model contracts","requirement_strength":"not specified","issue":"Neither HTTP nor webhook conventions define fan-out concurrency, result ordering, cancellation, or durability. One goroutine per endpoint is an unbounded resource risk.","interpretations":["Spawn freely, run serially, return completion order, stop on first failure, or use a fixed worker set."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"`FanOut` rejects input beyond `MaxFanOut`, requires a positive worker limit, uses a fixed worker bound, and returns one result per input in input order. Cancellation prevents new useful work but does not claim to interrupt a remote side effect already in flight. Fan-out is not a durable queue.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestFanOutBoundsConcurrencyAndPreservesResultOrder","TestFanOutRejectsUnboundedInputs"],"fixture_evidence":[],"fuzz_evidence":[],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["FanOut"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://raw.githubusercontent.com/golang/go/go1.26.6/doc/go_spec.html

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
implementation-defined behavior
application-policy
Go 1.26.6 cryptography, HTTP, URL, address, time, and encoding contracts
Go 1.26.6
go-language-source
https://raw.githubusercontent.com/golang/go/go1.26.6/doc/go_spec.html
Go language concurrency and memory-model contracts
not specified
Neither HTTP nor webhook conventions define fan-out concurrency, result ordering, cancellation, or durability. One goroutine per endpoint is an unbounded resource risk.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
`FanOut` rejects input beyond `MaxFanOut`, requires a positive worker limit, uses a fixed worker bound, and returns one result per input in input order. Cancellation prevents new useful work but does not claim to interrupt a remote side effect already in flight. Fan-out is not a durable queue.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Spawn freely, run serially, return completion order, stop on first failure, or use a fixed worker set.
TestFanOutBoundsConcurrencyAndPreservesResultOrder
TestFanOutRejectsUnboundedInputs
FanOut
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-023: Observation fields and trace propagation

**Authoritative reference:** [W3C Trace Context](https://www.w3.org/TR/2021/REC-trace-context-1-20211123/).

- **Status, owner, and classification:** `resolved`; maintainers; defensive
  privacy policy plus optional W3C Trace Context interoperability.
- **Source and issue:** W3C Trace Context defines propagation fields but not
  webhook signature coverage or application telemetry schemas. Logging raw
  payloads, signatures, URLs, IDs, or attacker strings can leak credentials and
  create high-cardinality telemetry.
- **Interpretations and peer behavior:** Sign trace fields, log full requests,
  attach raw errors, disable propagation, or expose a closed semantic schema.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  Core observations contain only fixed
  operation, outcome, reason, algorithm, status class, attempt, and bounded
  duration fields. They exclude payload, signature, secret, endpoint, event ID,
  replay key, and raw error text. The optional telemetry wrapper injects trace
  context after signing, so trace rotation does not invalidate the MAC, and
  preserves the caller's secure transport policy.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestVerifierObservesVerificationAndReplayWithoutSensitiveData`,
  `TestObserverWritesOnlyFixedSecretSafeAttributesThroughGoLog`, and
  `TestInstrumentHTTPClientInjectsTraceAndPreservesClientPolicy` cover
  `Observer` and adapters. Reconsider if a new signed trace-binding profile is
  explicitly required and its propagation lifecycle is defined.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-023","title":"Observation fields and trace propagation","status":"resolved","owner":"`webhook` maintainers","classification":"optional behavior","decision_scope":"defensive","specification":"W3C Trace Context Level 1","version":"W3C Recommendation 2021-11-23","source_authority":"trace-context-source","section":"W3C Trace Context Level 1 Section 3","requirement_strength":"not specified","issue":"W3C Trace Context defines propagation fields but not webhook signature coverage or application telemetry schemas. Logging raw payloads, signatures, URLs, IDs, or attacker strings can leak credentials and create high-cardinality telemetry.","interpretations":["Sign trace fields, log full requests, attach raw errors, disable propagation, or expose a closed semantic schema."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"Core observations contain only fixed operation, outcome, reason, algorithm, status class, attempt, and bounded duration fields. They exclude payload, signature, secret, endpoint, event ID, replay key, and raw error text. The optional telemetry wrapper injects trace context after signing, so trace rotation does not invalidate the MAC, and preserves the caller's secure transport policy.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestInstrumentHTTPClientInjectsTraceAndPreservesClientPolicy","TestVerifierObservesVerificationAndReplayWithoutSensitiveData","TestObserverRecordsBoundedMetricsAndCurrentSpanEvent"],"fixture_evidence":[],"fuzz_evidence":[],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["Observer","Observation","InstrumentHTTPClient"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.w3.org/TR/2021/REC-trace-context-1-20211123/

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
optional behavior
defensive
W3C Trace Context Level 1
W3C Recommendation 2021-11-23
trace-context-source
https://www.w3.org/TR/2021/REC-trace-context-1-20211123/
W3C Trace Context Level 1 Section 3
not specified
W3C Trace Context defines propagation fields but not webhook signature coverage or application telemetry schemas. Logging raw payloads, signatures, URLs, IDs, or attacker strings can leak credentials and create high-cardinality telemetry.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
Core observations contain only fixed operation, outcome, reason, algorithm, status class, attempt, and bounded duration fields. They exclude payload, signature, secret, endpoint, event ID, replay key, and raw error text. The optional telemetry wrapper injects trace context after signing, so trace rotation does not invalidate the MAC, and preserves the caller's secure transport policy.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Sign trace fields, log full requests, attach raw errors, disable propagation, or expose a closed semantic schema.
TestInstrumentHTTPClientInjectsTraceAndPreservesClientPolicy
TestVerifierObservesVerificationAndReplayWithoutSensitiveData
TestObserverRecordsBoundedMetricsAndCurrentSpanEvent
Observer
Observation
InstrumentHTTPClient
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-024: Provider presets remain unsupported

**Authoritative reference:** [RFC 2104](https://www.rfc-editor.org/rfc/rfc2104.html).

- **Status, owner, and classification:** `resolved`; maintainers; conformance
  claim and extension policy.
- **Source and issue:** Vendor webhook schemes differ in secret encoding,
  timestamp grammar, canonical bytes, multiple signatures, rotation, replay,
  and provider retry behavior. Similar use of HMAC is not interoperability.
- **Interpretations and peer behavior:** Market generic HMAC as compatible,
  embed provider switches in core, copy SDK snippets, or require isolated
  authoritative profiles and vectors.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  No provider preset is supported.
  Vendor names cannot enter the supported matrix until an isolated package has
  authoritative versioned documentation, independent positive vectors,
  negative mutations, rotation and retry semantics, and a maintenance owner.
  Generic `v1` remains provider-independent.
- **Evidence, public surface, upstream, and reconsideration:**
  `docs/providers.md`, `TestIndependentInteroperabilityVectors`, and absence of
  provider production packages are the current evidence. There is no upstream
  issue. Reconsider one provider at a time when complete authoritative evidence
  and ownership exist.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-024","title":"Provider presets remain unsupported","status":"resolved","owner":"`webhook` maintainers","classification":"interoperability policy","decision_scope":"extension-specific","specification":"RFC 2104 HMAC and RFC 4231 HMAC-SHA-256/HMAC-SHA-512 vectors","version":"RFC 2104","source_authority":"rfc2104-source","section":"RFC 2104 Sections 2-3","requirement_strength":"not specified","issue":"Vendor webhook schemes differ in secret encoding, timestamp grammar, canonical bytes, multiple signatures, rotation, replay, and provider retry behavior. Similar use of HMAC is not interoperability.","interpretations":["Market generic HMAC as compatible, embed provider switches in core, copy SDK snippets, or require isolated authoritative profiles and vectors."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"No provider preset is supported. Vendor names cannot enter the supported matrix until an isolated package has authoritative versioned documentation, independent positive vectors, negative mutations, rotation and retry semantics, and a maintenance owner. Generic `v1` remains provider-independent.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestIndependentInteroperabilityVectors","TestParseSignatureHeadersRejectsMalformedOrAmbiguousInput"],"fixture_evidence":[],"fuzz_evidence":[],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["Signer","Verifier"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc2104.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
interoperability policy
extension-specific
RFC 2104 HMAC and RFC 4231 HMAC-SHA-256/HMAC-SHA-512 vectors
RFC 2104
rfc2104-source
https://www.rfc-editor.org/rfc/rfc2104.txt
RFC 2104 Sections 2-3
not specified
Vendor webhook schemes differ in secret encoding, timestamp grammar, canonical bytes, multiple signatures, rotation, replay, and provider retry behavior. Similar use of HMAC is not interoperability.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
No provider preset is supported. Vendor names cannot enter the supported matrix until an isolated package has authoritative versioned documentation, independent positive vectors, negative mutations, rotation and retry semantics, and a maintenance owner. Generic `v1` remains provider-independent.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Market generic HMAC as compatible, embed provider switches in core, copy SDK snippets, or require isolated authoritative profiles and vectors.
TestIndependentInteroperabilityVectors
TestParseSignatureHeadersRejectsMalformedOrAmbiguousInput
Signer
Verifier
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-025: Delivery queue wire encoding

**Authoritative reference:** [RFC 8259](https://www.rfc-editor.org/rfc/rfc8259.html).

- **Status, owner, and classification:** `resolved`; maintainers; local JSON
  persistence boundary based on RFC 8259 and Go encoding contracts.
- **Source and issue:** RFC 8259 defines JSON but not webhook delivery fields,
  URL encoding, unknown-member policy, deterministic output, or size limits for
  queue/outbox transport.
- **Interpretations and peer behavior:** Serialize the public Go struct
  directly, use gob, preserve arbitrary headers, accept unknown fields, or
  define a private bounded versioned shape.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  Marshal a dedicated `v1` JSON shape
  with endpoint string, copied body, IDs, headers, and metadata under an exact
  byte limit. Decode one JSON value with unknown-field rejection and no trailing
  data, reconstruct and validate the endpoint, copy mutable collections, and
  reject unsafe required fields. Deterministic Go JSON output is contractual
  for equivalent values.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestDeliveryRequestWireRoundTripIsDeterministic`,
  `TestDeliveryRequestWireRejectsLimitsAndMalformedData`,
  `TestDeliveryWireRejectsEachUnsafeFieldIndependently`, and
  `FuzzDeliveryWire` cover `MarshalDeliveryRequest` and
  `UnmarshalDeliveryRequest`. Reconsider with a new wire version and migration
  plan before changing persisted bytes.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-025","title":"Delivery queue wire encoding","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"application-policy","specification":"RFC 8259 JSON","version":"RFC 8259","source_authority":"rfc8259-source","section":"RFC 8259 Sections 4 and 9","requirement_strength":"not specified","issue":"RFC 8259 defines JSON but not webhook delivery fields, URL encoding, unknown-member policy, deterministic output, or size limits for queue/outbox transport.","interpretations":["Serialize the public Go struct directly, use gob, preserve arbitrary headers, accept unknown fields, or define a private bounded versioned shape."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"Marshal a dedicated `v1` JSON shape with endpoint string, copied body, IDs, headers, and metadata under an exact byte limit. Decode one JSON value with unknown-field rejection and no trailing data, reconstruct and validate the endpoint, copy mutable collections, and reject unsafe required fields. Deterministic Go JSON output is contractual for equivalent values.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestDeliveryRequestWireRoundTripIsDeterministic","TestDeliveryRequestWireRejectsLimitsAndMalformedData","TestDeliveryWireRejectsEachUnsafeFieldIndependently"],"fixture_evidence":[],"fuzz_evidence":["FuzzDeliveryWire"],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["MarshalDeliveryRequest","UnmarshalDeliveryRequest"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc8259.txt

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
application-policy
RFC 8259 JSON
RFC 8259
rfc8259-source
https://www.rfc-editor.org/rfc/rfc8259.txt
RFC 8259 Sections 4 and 9
not specified
RFC 8259 defines JSON but not webhook delivery fields, URL encoding, unknown-member policy, deterministic output, or size limits for queue/outbox transport.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
Marshal a dedicated `v1` JSON shape with endpoint string, copied body, IDs, headers, and metadata under an exact byte limit. Decode one JSON value with unknown-field rejection and no trailing data, reconstruct and validate the endpoint, copy mutable collections, and reject unsafe required fields. Deterministic Go JSON output is contractual for equivalent values.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Serialize the public Go struct directly, use gob, preserve arbitrary headers, accept unknown fields, or define a private bounded versioned shape.
TestDeliveryRequestWireRoundTripIsDeterministic
TestDeliveryRequestWireRejectsLimitsAndMalformedData
TestDeliveryWireRejectsEachUnsafeFieldIndependently
FuzzDeliveryWire
MarshalDeliveryRequest
UnmarshalDeliveryRequest
docs/specification-decisions.md
</pre>

</details>

## WEBHOOK-DEC-026: Compatibility and reconsideration policy

**Authoritative reference:** [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html).

- **Status, owner, and classification:** `resolved`; maintainers; SemVer and
  change-control policy.
- **Source and issue:** None of the external specifications determines which
  local webhook choices are public compatibility surfaces or how security
  tightening is communicated.
- **Interpretations and peer behavior:** Treat only Go identifiers as API,
  silently tighten parsers, evolve canonical bytes in place, or version every
  observable protocol and operational classification.
- **Selected behavior and consequences:**
  Security, resource, compatibility, and wire consequences are included in
  the selected behavior below.
  Canonical bytes, field grammar,
  algorithm identifiers, envelope and delivery JSON, exported error identity,
  retry classification, replay identity, endpoint defaults, and established
  provider profiles are compatibility surfaces. Wire changes require a new
  protocol version and normally major impact. Security fixes may reject
  formerly accepted unsafe input but require explicit changelog and migration
  review. Unknown or unresolved behavior cannot become a silent default.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestIndependentInteroperabilityVectors`,
  `TestDeliveryRequestWireRoundTripIsDeterministic`, API baselines, golden
  vectors, `docs/migration.md`, and conformance checks enforce the current
  surfaces. Reconsider when a superseding decision records source, migration,
  executable evidence, and release impact; preserve this entry as history.


<details>
<summary>Machine-auditable bindings</summary>

```json
{"id":"WEBHOOK-DEC-026","title":"Compatibility and reconsideration policy","status":"resolved","owner":"`webhook` maintainers","classification":"omission","decision_scope":"application-policy","specification":"Semantic Versioning 2.0.0","version":"Semantic Versioning 2.0.0","source_authority":"semver-source","section":"Semantic Versioning 2.0.0 clauses 1-8","requirement_strength":"not specified","issue":"None of the external specifications determines which local webhook choices are public compatibility surfaces or how security tightening is communicated.","interpretations":["Treat only Go identifiers as API, silently tighten parsers, evolve canonical bytes in place, or version every observable protocol and operational classification."],"peer_behavior":"Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.","selected_behavior":"Canonical bytes, field grammar, algorithm identifiers, envelope and delivery JSON, exported error identity, retry classification, replay identity, endpoint defaults, and established provider profiles are compatibility surfaces. Wire changes require a new protocol version and normally major impact. Security fixes may reject formerly accepted unsafe input but require explicit changelog and migration review. Unknown or unresolved behavior cannot become a silent default.","rationale":"The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.","security_consequences":"The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.","resource_consequences":"The selected behavior retains the documented finite processing, allocation, and ownership bounds.","compatibility_consequences":"Changing this selection requires compatibility and specification-decision review.","wire_consequences":"Wire-visible output and rejection behavior follow the selected behavior exactly.","executable_evidence":["TestIndependentInteroperabilityVectors","TestDeliveryRequestWireRoundTripIsDeterministic"],"fixture_evidence":[],"fuzz_evidence":[],"interoperability_evidence":[],"differential_evidence":[],"public_apis":["Canonicalize","Signature","Envelope","DeliveryRequest","VerificationError"],"documentation":["docs/specification-decisions.md"],"upstream_status":"No upstream issue or erratum is currently recorded for this decision.","reconsider_when":"A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary."}
```

Authority URL: https://semver.org/spec/v2.0.0.html

Exact-value mirror:

<pre>
resolved
`webhook` maintainers
omission
application-policy
Semantic Versioning 2.0.0
Semantic Versioning 2.0.0
semver-source
https://semver.org/spec/v2.0.0.html
Semantic Versioning 2.0.0 clauses 1-8
not specified
None of the external specifications determines which local webhook choices are public compatibility surfaces or how security tightening is communicated.
Maintained-peer behavior has not been assessed for this decision; the current record relies on the cited authorities and local executable evidence.
Canonical bytes, field grammar, algorithm identifiers, envelope and delivery JSON, exported error identity, retry classification, replay identity, endpoint defaults, and established provider profiles are compatibility surfaces. Wire changes require a new protocol version and normally major impact. Security fixes may reject formerly accepted unsafe input but require explicit changelog and migration review. Unknown or unresolved behavior cannot become a silent default.
The choice makes the observable boundary explicit and keeps it reviewable against the cited authority and executable evidence.
The selected behavior prevents ambiguous input or lifecycle state from silently widening the accepted policy.
The selected behavior retains the documented finite processing, allocation, and ownership bounds.
Changing this selection requires compatibility and specification-decision review.
Wire-visible output and rejection behavior follow the selected behavior exactly.
No upstream issue or erratum is currently recorded for this decision.
A monitored authority, maintained peer, security review, or supported Go contract changes this observable boundary.
Treat only Go identifiers as API, silently tighten parsers, evolve canonical bytes in place, or version every observable protocol and operational classification.
TestIndependentInteroperabilityVectors
TestDeliveryRequestWireRoundTripIsDeterministic
Canonicalize
Signature
Envelope
DeliveryRequest
VerificationError
docs/specification-decisions.md
</pre>

</details>

## Unresolved decisions

None for the currently supported provider-independent webhook surfaces. New
provider profiles, algorithms, canonicalization rules, transport mappings,
errata, or peer divergences MUST be registered before observable behavior is
selected. An unresolved wire, security, resource, or lifecycle decision blocks
the affected release claim.

[RFC2119]: https://www.rfc-editor.org/rfc/rfc2119
[RFC8174]: https://www.rfc-editor.org/rfc/rfc8174
