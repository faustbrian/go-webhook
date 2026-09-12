# Specification provenance

`manifest.tsv` pins every external authority used by the local webhook protocol
and its HTTP, replay, delivery, envelope, and endpoint-security boundaries. The
[specification decision register](../docs/specification-decisions.md) records
selected behavior without claiming RFC 9421, CloudEvents, or provider-profile
compatibility.

## Decision conformance matrix

| Decision | Primary authority | Executable evidence | Differential evidence |
| --- | --- | --- | --- |
| WEBHOOK-DEC-001 | RFC 9421 HTTP Message Signatures comparison boundary | `TestCanonicalizeProducesStableVersionedBytes`, `TestParseSignatureHeadersRejectsMalformedOrAmbiguousInput`, `TestIndependentInteroperabilityVectors` | Not assessed |
| WEBHOOK-DEC-002 | RFC 2104 HMAC and RFC 4231 HMAC-SHA-256/HMAC-SHA-512 vectors | `TestSignerAndVerifierSupportSHA256AndSHA512`, `TestVerifierRejectsMutationOfEverySignedComponent`, `TestIndependentInteroperabilityVectors` | Not assessed |
| WEBHOOK-DEC-003 | RFC 4648 Base-N Encodings | `TestCanonicalizeProducesStableVersionedBytes`, `TestIndependentInteroperabilityVectors` | Not assessed |
| WEBHOOK-DEC-004 | RFC 3986 URI Generic Syntax | `TestVerifierBindsDuplicateQueryValueOrder`, `TestVerifierRejectsMutationOfEverySignedComponent`, `TestSignAndVerifyRequestUsesRawBodyAndRestoresIt`, `TestIndependentInteroperabilityVectors` | Not assessed |
| WEBHOOK-DEC-005 | RFC 9110 HTTP Semantics | `TestVerifyRequestRejectsMutationOfFixedSignedHeaders`, `TestVerifyRequestRejectsDuplicateSignedHeaderBeforeBodyRead` | Not assessed |
| WEBHOOK-DEC-006 | RFC 3339 Internet date and time | `TestVerifierComparesCallerTimestampAtProtocolSecondPrecision`, `TestVerifierTimestampToleranceBoundaries` | Not assessed |
| WEBHOOK-DEC-007 | RFC 2104 HMAC and RFC 4231 HMAC-SHA-256/HMAC-SHA-512 vectors | `TestSignerUsesOneInjectedNonceAcrossRotationSignatures`, `TestSignerSelectsRotationKeyAtSignedTimestamp`, `TestRotationSignsAllActiveKeysAndAcceptsOverlap` | Not assessed |
| WEBHOOK-DEC-008 | RFC 8941 Structured Fields comparison boundary | `TestParseSignatureHeadersRejectsMalformedOrAmbiguousInput`, `TestParseSignatureHeadersAppliesLimitsBeforeDecoding`, `TestSignatureHeadersRoundTripMultipleRotationSignatures` | Not assessed |
| WEBHOOK-DEC-009 | RFC 9110 HTTP Semantics | `TestCaptureBodyPreservesExactBytesAndRestoresRequest`, `TestCaptureBodyBoundsUnknownLengthBeforeAllocation`, `TestCaptureBodyAfterPriorReadAuthenticatesOnlyRemainingBytes` | Not assessed |
| WEBHOOK-DEC-010 | RFC 2104 HMAC and RFC 4231 HMAC-SHA-256/HMAC-SHA-512 vectors | `TestVerifierSkipsMalformedAndInactiveSignaturesDeterministically`, `TestVerificationErrorMethodsAreNilSafe` | Not assessed |
| WEBHOOK-DEC-011 | RFC 9110 HTTP Semantics | `TestVerifyAndRecordAtomicallyRejectsReplay`, `TestReplayIdentitySurvivesSecretRotation`, `TestVerifyAndRecordFailsClosedForMissingIDAndStoreFailure` | Not assessed |
| WEBHOOK-DEC-012 | RFC 9110 HTTP Semantics | `TestVerifyRequestExtractsEventIDAfterAuthentication`, `TestHeaderEventIDRejectsDuplicateAndOversizedValues` | Not assessed |
| WEBHOOK-DEC-013 | CloudEvents 1.0.2 comparison boundary | `TestEnvelopeMarshalIsDeterministicAndPreservesData`, `TestEnvelopeRejectsInvalidRequiredFieldsAndData` | Not assessed |
| WEBHOOK-DEC-014 | RFC 9110 HTTP Semantics | `TestDeliverPreservesRequestStatusAndResponseBoundaries`, `TestVerifyRequestRejectsMutationOfFixedSignedHeaders` | Not assessed |
| WEBHOOK-DEC-015 | RFC 9110 HTTP Semantics | `TestDeliverRetriesRetryableStatusAndHonorsRetryAfter`, `TestRetryAfterSupportsHTTPDateAndCapsDelay`, `TestRetryPolicyCapsOverflowingRetryAfterSeconds` | Not assessed |
| WEBHOOK-DEC-016 | RFC 9110 HTTP Semantics | `TestDeliverRetriesNetworkFailureAndStopsAtBound`, `TestDeliverOnceDisablesInternalRetries`, `TestHandleUsesSingleDeliveryAttempt`, `TestPublisherPerformsSingleAttemptForRelay` | Not assessed |
| WEBHOOK-DEC-017 | RFC 9110 HTTP Semantics | `TestDeliverBoundsResponseBody`, `TestReadResponseRejectsMissingReadAndCloseFailures`, `TestDeliverClassifiesCanceledTransportAndExhaustedStatus` | Not assessed |
| WEBHOOK-DEC-018 | RFC 9110 HTTP Semantics | `TestDeliverDoesNotRetryTerminalStatusAndDeadLetters`, `TestDeliveryPreservesDeadLetterHookFailure`, `TestReplayAuditsBeforeStartingNewDelivery` | Not assessed |
| WEBHOOK-DEC-019 | RFC 3986 URI Generic Syntax | `TestSSRFPolicyRejectsUnsafeURLsAndAddresses`, `TestSSRFPolicyAllowsExplicitPrefixOnlyWhenConfigured` | Not assessed |
| WEBHOOK-DEC-020 | IANA IPv4 and IPv6 Special-Purpose Address Registries at 2025-10-09 | `TestSSRFPolicyRejectsUnsafeURLsAndAddresses`, `TestSSRFPolicyRejectsMixedAndOversizedDNSAnswers` | Not assessed |
| WEBHOOK-DEC-021 | Go 1.26.6 cryptography, HTTP, URL, address, time, and encoding contracts | `TestSecureHTTPClientRevalidatesDNSAtDialTime`, `TestSecureHTTPClientRejectsRedirectWithoutContactingTarget`, `TestSecureHTTPClientConfigurationAndTransportFailures` | Not assessed |
| WEBHOOK-DEC-022 | Go 1.26.6 cryptography, HTTP, URL, address, time, and encoding contracts | `TestFanOutBoundsConcurrencyAndPreservesResultOrder`, `TestFanOutRejectsUnboundedInputs` | Not assessed |
| WEBHOOK-DEC-023 | W3C Trace Context Level 1 | `TestInstrumentHTTPClientInjectsTraceAndPreservesClientPolicy`, `TestVerifierObservesVerificationAndReplayWithoutSensitiveData`, `TestObserverRecordsBoundedMetricsAndCurrentSpanEvent` | Not assessed |
| WEBHOOK-DEC-024 | RFC 2104 HMAC and RFC 4231 HMAC-SHA-256/HMAC-SHA-512 vectors | `TestIndependentInteroperabilityVectors`, `TestParseSignatureHeadersRejectsMalformedOrAmbiguousInput` | Not assessed |
| WEBHOOK-DEC-025 | RFC 8259 JSON | `TestDeliveryRequestWireRoundTripIsDeterministic`, `TestDeliveryRequestWireRejectsLimitsAndMalformedData`, `TestDeliveryWireRejectsEachUnsafeFieldIndependently` | Not assessed |
| WEBHOOK-DEC-026 | Semantic Versioning 2.0.0 | `TestIndependentInteroperabilityVectors`, `TestDeliveryRequestWireRoundTripIsDeterministic` | Not assessed |

Maintained-provider interoperability and maintained-peer differential behavior
are currently not assessed. The first four decisions retain reproducible
independent Python standard-library fixture evidence without classifying that
primitive comparison as maintained-peer agreement.

Run `golib specification check` for structural validation and
`golib specification check --online` to revalidate monitored authority bytes.
Reviewed changes to monitored authorities are preserved in the append-only
[upstream authority review history](upstream-reviews.md).
