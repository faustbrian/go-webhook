# Signature scheme v1

`v1` supports `sha256` and `sha512`. It signs a labeled, newline-delimited
sequence of: version, algorithm, Unix timestamp seconds, nonce, key ID,
uppercase method, escaped path, canonical query, lowercase host, SHA-256 body
digest, and canonical metadata. Variable byte fields are unpadded base64url,
so embedded delimiters cannot alter the grammar. Query keys and values use Go
URL parsing and sorting.
Metadata keys are sorted; key and value bytes are independently base64url
encoded without padding.

The HTTP header is a structured, single-value-per-line field:

```text
Webhook-Signature: v1;algorithm=sha256;keyid=<base64url>;timestamp=<unix>;nonce=<base64url>;signature=<base64url>
```

Unknown versions or algorithms, duplicate parameters, padding, whitespace,
invalid encodings, timestamp disagreement, and duplicate complete signature
values are rejected. Comparisons use `hmac.Equal`. The timestamp is covered by
the MAC and checked against the verifier-owned clock and inclusive tolerance.
The nonce is also covered by the MAC, bounded to 128 UTF-8 bytes, and generated
once per signing operation. `SignerConfig.NonceGenerator` supports
deterministic tests; the default uses `crypto/rand`.

The exact normative vectors are `testdata/vectors/v1.json`. They were produced
independently by Python's standard `hmac`, `hashlib`, and URL primitives and
are checked by `scripts/check_interoperability.py`. Any canonicalization or
header change requires a major version and new vector version.
