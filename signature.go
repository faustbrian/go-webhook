// Package webhook provides protocol-independent webhook signing, verification,
// replay protection, and delivery primitives. Start with NewSigner/NewVerifier
// and their HTTP helpers; NewDeliverer composes caller-owned HTTP and endpoint
// policy. Local v1 is not RFC 9421, CloudEvents or a vendor preset. Optional
// adapters import existing queue, outbox, idempotency, logging and telemetry
// owners; the core owns none of their lifecycles and introduces no cgo.
//
// Configuration key material is copied, not erased; message slices/maps and
// callbacks are borrowed. Bound direct inputs and make shared collaborators
// concurrency-safe. No clock, callback or external I/O is forcibly preempted.
// Verification authenticates exact remaining body bytes, not application
// authorization or exactly-once business effects. Secure egress requires both
// endpoint policy and a transport enforcing validated dial targets.
package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/url"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

// Algorithm identifies a supported HMAC construction.
type Algorithm string

const (
	// SHA256 is HMAC-SHA-256.
	SHA256 Algorithm = "sha256"
	// SHA512 is HMAC-SHA-512.
	SHA512 Algorithm = "sha512"
)

var (
	// ErrInvalidSignature means that no supplied signature authenticated.
	ErrInvalidSignature = errors.New("invalid webhook signature")
	// ErrInvalidTimestamp means that a signature timestamp is missing, invalid,
	// or outside the configured tolerance.
	ErrInvalidTimestamp = errors.New("invalid webhook timestamp")
	// ErrInvalidConfiguration means a signer or verifier configuration is not
	// safe to use.
	ErrInvalidConfiguration = errors.New("invalid webhook configuration")
	// ErrNoActiveKey means no configured signing key is active.
	ErrNoActiveKey = errors.New("no active webhook signing key")
	// ErrReplay means an authenticated event ID was already atomically stored.
	ErrReplay = errors.New("webhook replay detected")
	// ErrMissingEventID means replay protection was requested without an ID.
	ErrMissingEventID = errors.New("webhook event ID is required")
	// ErrReplayStore means replay state could not be checked and the request was
	// rejected closed.
	ErrReplayStore = errors.New("webhook replay store unavailable")
	// ErrNonceGeneration means a signer could not produce a safe nonce.
	ErrNonceGeneration = errors.New("webhook nonce generation failed")
)

const maxNonceBytes = 128

// NonceGenerator returns a fresh public nonce for one signing operation.
// Implementations must be concurrency-safe and bound their work. Return fresh
// nonempty UTF-8 of at most 128 bytes; errors/invalid outputs become
// ErrNonceGeneration. Signing contains no callback panic or timeout boundary.
type NonceGenerator func() (string, error)

// ReplayStore atomically checks and records replay keys. Implementations MUST
// return true only when the key was absent and was recorded with expiresAt in
// the same atomic operation. false, nil means the key already existed.
type ReplayStore interface {
	// CheckAndRecord must atomically record absent identity until expiresAt;
	// retained identity returns false. Honor ctx and retain no unbounded
	// application data. Core keys are domain/namespace/event-ID digests. Storage
	// errors fail closed; admission is not atomic with downstream business effects.
	CheckAndRecord(ctx context.Context, key string, expiresAt time.Time) (recorded bool, err error)
}

// Message contains the exact request components covered by a signature. Body
// is never transformed before hashing. Slices/maps are borrowed for a call;
// concurrent mutation is forbidden. Direct APIs impose no overall body/query/
// metadata/key-ID bound; the caller must cap untrusted inputs before allocating.
type Message struct {
	// Timestamp is signed at nonnegative Unix-second precision; zero lets Sign sample its clock.
	Timestamp time.Time
	// Nonce is signed public UTF-8, 1–128 bytes; empty lets Sign generate one.
	Nonce string
	// Method is nonempty case-sensitive method bytes with no CR/LF.
	Method string
	// Path is caller-supplied path bytes; HTTP helpers use EscapedPath or /.
	Path string
	// RawQuery is parsed with Go form semantics; key order normalizes, duplicate value order stays.
	RawQuery string
	// Host is lowercased in the canonical form, including any explicit port.
	Host string
	// ContentType is exact signed bytes; HTTP helpers admit one value up to 256 UTF-8 bytes.
	ContentType string
	// IdempotencyKey is exact signed identity, not proof of receiver deduplication.
	IdempotencyKey string
	// Body is borrowed exact bytes hashed with SHA-256 even for HMAC-SHA-512.
	Body []byte
	// Metadata is borrowed signed string data sorted by key; nil means no entries.
	Metadata map[string]string
}

// SigningKey is a versioned outbound secret. Empty validity bounds are open.
type SigningKey struct {
	// ID is a required unique configured key identity, not secret material.
	ID string
	// Secret is required nonempty HMAC material copied by NewSigner; entropy is caller-owned.
	Secret []byte
	// NotBefore is the inclusive start at signed time; zero leaves the start open.
	NotBefore time.Time
	// NotAfter is the inclusive end at signed time; zero leaves the end open.
	NotAfter time.Time
	// Revoked excludes the key regardless of its time window.
	Revoked bool
}

// VerificationKey is a versioned inbound secret. Empty validity bounds are
// open.
type VerificationKey struct {
	// ID is a required unique identity selected from the candidate signature.
	ID string
	// Secret is required nonempty HMAC material copied by NewVerifier; entropy is caller-owned.
	Secret []byte
	// NotBefore is the inclusive validity start at the candidate signed time; zero is open.
	NotBefore time.Time
	// NotAfter is the inclusive validity end at the candidate signed time; zero is open.
	NotAfter time.Time
	// Revoked rejects this key regardless of its time window.
	Revoked bool
}

// Signature is a transport-neutral signature value.
type Signature struct {
	// Version identifies local v1, not an RFC 9421 or vendor format.
	Version string
	// Algorithm is the configured HMAC construction; inbound values cannot widen it.
	Algorithm Algorithm
	// KeyID attributes the signing key and is included in canonical bytes.
	KeyID string
	// Timestamp is normalized signed UTC Unix-second time.
	Timestamp time.Time
	// Nonce is the signed public nonce shared by all keys in one Sign call.
	Nonce string
	// Value is an unpadded base64url MAC, not ciphertext or secret key material.
	Value string
}

// Verification describes the key and algorithm that authenticated a message.
type Verification struct {
	// KeyID identifies the first candidate that authenticated.
	KeyID string
	// Algorithm is the configured accepted HMAC construction.
	Algorithm Algorithm
	// Timestamp is the accepted signature time, not processing completion time.
	Timestamp time.Time
	// Nonce is the accepted signed nonce; Verify alone does not store it.
	Nonce string
}

// VerificationError separates a stable caller-facing message from internal
// diagnostics. Kind may wrap body-I/O diagnostics; exported caller-created
// fields are not sanitized. Use SafeMessage for untrusted HTTP responses; Error
// renders Kind.Error verbatim, not a redaction guarantee. Methods accept nil.
type VerificationError struct {
	// Kind is the wrapped category; caller-constructed errors render its Error verbatim.
	Kind error
	// Diagnostic is internal text not included by Error or SafeMessage; caller values are not sanitized.
	Diagnostic string
}

// Error implements error without exposing internal diagnostic details.
func (e *VerificationError) Error() string {
	if e == nil || e.Kind == nil {
		return "webhook verification failed"
	}

	return e.Kind.Error()
}

// Unwrap enables errors.Is checks against the failure category.
func (e *VerificationError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Kind
}

// SafeMessage returns a stable response suitable for an untrusted caller.
func (e *VerificationError) SafeMessage() string {
	return "webhook verification failed"
}

// Canonicalize constructs owned local v1 bytes with a fixed label order/final
// LF, raw base64url variable fields, sorted metadata and a SHA-256 body digest.
// Query uses Go ParseQuery/Values.Encode, retaining duplicate value order; host
// lowercases and method case stays exact. This is local policy under
// WEBHOOK-DEC-001–005 in docs/specification-decisions.md, protected by
// TestCanonicalizeProducesStableVersionedBytes and independent v1 vectors.
// It neither validates endpoint safety nor authorizes a request.
func Canonicalize(message Message, keyID string, algorithm Algorithm) ([]byte, error) {
	if _, err := hashFactory(algorithm); err != nil {
		return nil, err
	}
	if message.Timestamp.IsZero() {
		return nil, fmt.Errorf("%w: timestamp is required", ErrInvalidTimestamp)
	}
	if message.Timestamp.Unix() < 0 {
		return nil, fmt.Errorf("%w: timestamp is required", ErrInvalidTimestamp)
	}
	if !validNonce(message.Nonce) {
		return nil, fmt.Errorf("%w: nonce is required and must be bounded UTF-8", ErrInvalidConfiguration)
	}
	if message.Method == "" || strings.ContainsAny(message.Method, "\r\n") {
		return nil, fmt.Errorf("%w: method is required and cannot contain a line break", ErrInvalidConfiguration)
	}

	query, err := canonicalQuery(message.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid query: %v", ErrInvalidConfiguration, err)
	}

	bodyDigest := sha256.Sum256(message.Body)
	metadata := canonicalMetadata(message.Metadata)
	encode := base64.RawURLEncoding.EncodeToString
	canonical := strings.Join([]string{
		"webhook-v1",
		"algorithm:" + string(algorithm),
		fmt.Sprintf("timestamp:%d", message.Timestamp.Unix()),
		"nonce:" + encode([]byte(message.Nonce)),
		"key-id:" + encode([]byte(keyID)),
		"method:" + message.Method,
		"path:" + encode([]byte(message.Path)),
		"query:" + encode([]byte(query)),
		"host:" + encode([]byte(strings.ToLower(message.Host))),
		"content-type:" + encode([]byte(message.ContentType)),
		"idempotency-key:" + encode([]byte(message.IdempotencyKey)),
		"body-sha256:" + encode(bodyDigest[:]),
		"metadata:" + encode([]byte(metadata)),
		"",
	}, "\n")

	return []byte(canonical), nil
}

func canonicalQuery(raw string) (string, error) {
	values, err := url.ParseQuery(raw)
	if err != nil {
		return "", err
	}

	return values.Encode(), nil
}

func canonicalMetadata(metadata map[string]string) string {
	keys := make([]string, 0, len(metadata))
	for key := range metadata {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	encode := base64.RawURLEncoding.EncodeToString
	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		lines = append(lines, encode([]byte(key))+"="+encode([]byte(metadata[key])))
	}

	return strings.Join(lines, "\n")
}

// SignerConfig configures outbound signatures.
type SignerConfig struct {
	// Algorithm must explicitly be SHA256 or SHA512; zero has no default.
	Algorithm Algorithm
	// Keys must be nonempty with unique nonempty IDs/secrets and non-inverted windows.
	Keys []SigningKey
	// Clock is borrowed; nil uses time.Now. Bound callbacks and make them concurrency-safe.
	Clock func() time.Time
	// NonceGenerator is borrowed; nil uses 18 crypto/rand bytes encoded as base64url.
	NonceGenerator NonceGenerator
}

// Signer creates signatures for every active key, newest first. Use NewSigner;
// nil/zero Signer is unusable. It is immutable and concurrently usable when
// borrowed Clock/NonceGenerator and input objects support that usage. No worker,
// key erasure or hot-reload mechanism is owned; reconstruct to change keys.
type Signer struct {
	algorithm Algorithm
	hash      func() hash.Hash
	keys      []SigningKey
	clock     func() time.Time
	nonce     NonceGenerator
}

// NewSigner validates key identity/windows and copies slice/secret bytes before
// publication, returning nil on error. Clock/nonce closures remain borrowed;
// no entropy-length minimum or upper bound on configured key count is imposed.
func NewSigner(config SignerConfig) (*Signer, error) {
	factory, err := hashFactory(config.Algorithm)
	if err != nil {
		return nil, err
	}
	if len(config.Keys) == 0 {
		return nil, fmt.Errorf("%w: at least one key is required", ErrInvalidConfiguration)
	}
	clock := config.Clock
	if clock == nil {
		clock = time.Now
	}
	nonce := config.NonceGenerator
	if nonce == nil {
		nonce = defaultNonce
	}
	keys := make([]SigningKey, len(config.Keys))
	keyIDs := make(map[string]struct{}, len(config.Keys))
	for index, key := range config.Keys {
		if key.ID == "" {
			return nil, fmt.Errorf("%w: key ID is required", ErrInvalidConfiguration)
		}
		if len(key.Secret) == 0 {
			return nil, fmt.Errorf("%w: key secret is required", ErrInvalidConfiguration)
		}
		if !key.NotBefore.IsZero() && !key.NotAfter.IsZero() && key.NotAfter.Before(key.NotBefore) {
			return nil, fmt.Errorf("%w: key validity window is invalid", ErrInvalidConfiguration)
		}
		if _, duplicate := keyIDs[key.ID]; duplicate {
			return nil, fmt.Errorf("%w: duplicate key ID", ErrInvalidConfiguration)
		}
		keyIDs[key.ID] = struct{}{}
		key.Secret = append([]byte(nil), key.Secret...)
		keys[index] = key
	}

	return &Signer{algorithm: config.Algorithm, hash: factory, keys: keys, clock: clock, nonce: nonce}, nil
}

// Sign creates one signature for each active key and returns owned values.
// Zero timestamp samples Clock, then truncates to UTC Unix seconds; one generated
// nonce is shared across keys without mutating Message. Key validity includes
// both endpoints at that signed time, not a separate current-time instant.
// No active key returns ErrNoActiveKey. Borrowed callbacks may block/panic.
func (s *Signer) Sign(message Message) ([]Signature, error) {
	now := s.clock().UTC()
	if message.Timestamp.IsZero() {
		message.Timestamp = now
	}
	message.Timestamp = time.Unix(message.Timestamp.Unix(), 0).UTC()
	if message.Nonce == "" {
		if s.nonce == nil {
			return nil, ErrNonceGeneration
		}
		nonce, err := s.nonce()
		if err != nil || !validNonce(nonce) {
			return nil, ErrNonceGeneration
		}
		message.Nonce = nonce
	} else if !validNonce(message.Nonce) {
		return nil, fmt.Errorf("%w: nonce is invalid", ErrInvalidConfiguration)
	}
	keys := make([]SigningKey, 0, len(s.keys))
	for _, key := range s.keys {
		if activeKey(key.NotBefore, key.NotAfter, key.Revoked, message.Timestamp) {
			keys = append(keys, key)
		}
	}
	if len(keys) == 0 {
		return nil, ErrNoActiveKey
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].NotBefore.Equal(keys[j].NotBefore) {
			return strings.Compare(keys[i].ID, keys[j].ID) == -1
		}

		return keys[i].NotBefore.After(keys[j].NotBefore)
	})

	signatures := make([]Signature, 0, len(keys))
	for _, key := range keys {
		canonical, err := Canonicalize(message, key.ID, s.algorithm)
		if err != nil {
			return nil, err
		}
		value := signWithHash(s.hash, key.Secret, canonical)
		signatures = append(signatures, Signature{
			Version:   "v1",
			Algorithm: s.algorithm,
			KeyID:     key.ID,
			Timestamp: message.Timestamp,
			Nonce:     message.Nonce,
			Value:     base64.RawURLEncoding.EncodeToString(value),
		})
	}

	return signatures, nil
}

// VerifierConfig configures inbound signature verification.
type VerifierConfig struct {
	// Algorithm must explicitly be SHA256 or SHA512; zero has no default.
	Algorithm Algorithm
	// Keys must be nonempty and valid; slice and secret bytes are copied.
	Keys []VerificationKey
	// Clock is borrowed; nil uses time.Now. Its time controls skew and replay expiry.
	Clock func() time.Time
	// Tolerance is inclusive absolute clock skew; zero requires the sampled instant, negative rejects.
	Tolerance time.Duration
	// ReplayStore is optional and borrowed; nil disables replay checks, typed nil is caller-invalid.
	ReplayStore ReplayStore
	// ReplayTTL must be positive with a store and zero without one; no maximum is imposed.
	ReplayTTL time.Duration
	// ReplayNamespace must be nonempty with a store and empty without one; caller defines collision scope.
	ReplayNamespace string
	// Observer is optional and borrowed; callbacks are synchronous and core invocation contains panics.
	Observer Observer
}

// Verifier authenticates signatures with configured clock skew. Use NewVerifier;
// nil/zero Verifier is unusable. Immutable copied keys support concurrent use
// when borrowed clock/store/observer and input data do. It owns no worker/store
// lifecycle and cannot couple replay admission to application transaction commit.
type Verifier struct {
	algorithm Algorithm
	hash      func() hash.Hash
	keys      map[string]VerificationKey
	clock     func() time.Time
	tolerance time.Duration
	replay    ReplayStore
	replayTTL time.Duration
	namespace string
	observer  Observer
}

// NewVerifier validates key/algorithm/skew/replay configuration, copies key
// material and returns nil on failure. Other collaborators are retained;
// typed-nil interfaces are not detected and clock panics are not contained.
func NewVerifier(config VerifierConfig) (*Verifier, error) {
	factory, err := hashFactory(config.Algorithm)
	if err != nil {
		return nil, err
	}
	if len(config.Keys) == 0 {
		return nil, fmt.Errorf("%w: keys are required", ErrInvalidConfiguration)
	}
	if config.Tolerance < 0 {
		return nil, fmt.Errorf("%w: keys are required and tolerance cannot be negative", ErrInvalidConfiguration)
	}
	if config.ReplayStore != nil {
		if config.ReplayTTL <= 0 {
			return nil, fmt.Errorf("%w: positive replay TTL is required", ErrInvalidConfiguration)
		}
		if config.ReplayNamespace == "" {
			return nil, fmt.Errorf("%w: replay namespace is required", ErrInvalidConfiguration)
		}
	} else {
		if config.ReplayTTL != 0 {
			return nil, fmt.Errorf("%w: replay TTL requires a store", ErrInvalidConfiguration)
		}
		if config.ReplayNamespace != "" {
			return nil, fmt.Errorf("%w: replay namespace requires a store", ErrInvalidConfiguration)
		}
	}
	clock := config.Clock
	if clock == nil {
		clock = time.Now
	}
	keys := make(map[string]VerificationKey, len(config.Keys))
	for _, key := range config.Keys {
		if key.ID == "" {
			return nil, fmt.Errorf("%w: key ID is required", ErrInvalidConfiguration)
		}
		if len(key.Secret) == 0 {
			return nil, fmt.Errorf("%w: key secret is required", ErrInvalidConfiguration)
		}
		if !key.NotBefore.IsZero() && !key.NotAfter.IsZero() && key.NotAfter.Before(key.NotBefore) {
			return nil, fmt.Errorf("%w: key validity window is invalid", ErrInvalidConfiguration)
		}
		if _, duplicate := keys[key.ID]; duplicate {
			return nil, fmt.Errorf("%w: duplicate key ID", ErrInvalidConfiguration)
		}
		key.Secret = append([]byte(nil), key.Secret...)
		keys[key.ID] = key
	}

	return &Verifier{
		algorithm: config.Algorithm,
		hash:      factory,
		keys:      keys,
		clock:     clock,
		tolerance: config.Tolerance,
		replay:    config.ReplayStore,
		replayTTL: config.ReplayTTL,
		namespace: config.ReplayNamespace,
		observer:  config.Observer,
	}, nil
}

// Verify accepts the first valid signature without exposing which checks
// failed. MAC equality uses hmac.Equal under local WEBHOOK-DEC-010 policy;
// metadata/key lookup/candidate ordering do not have constant total runtime.
// Input order selects the first valid key; invalid candidates do not poison
// later candidates. Zero Message timestamp/nonce adopts the candidate values.
// Verify performs no replay storage or observations and imposes no candidate
// count/encoded-value bound; use the bounded HTTP parser for untrusted headers.
func (v *Verifier) Verify(message Message, signatures []Signature) (Verification, error) {
	now := v.clock().UTC()
	for _, signature := range signatures {
		if signature.Version != "v1" || signature.Algorithm != v.algorithm || signature.Timestamp.IsZero() || !validNonce(signature.Nonce) {
			continue
		}
		if now.Sub(signature.Timestamp).Abs() > v.tolerance {
			continue
		}
		if !message.Timestamp.IsZero() && message.Timestamp.Unix() != signature.Timestamp.Unix() {
			continue
		}
		if message.Nonce != "" && message.Nonce != signature.Nonce {
			continue
		}
		key, exists := v.keys[signature.KeyID]
		if !exists || !activeKey(key.NotBefore, key.NotAfter, key.Revoked, signature.Timestamp) {
			continue
		}
		candidate, err := base64.RawURLEncoding.DecodeString(signature.Value)
		if err != nil {
			continue
		}
		candidateMessage := message
		candidateMessage.Timestamp = signature.Timestamp
		candidateMessage.Nonce = signature.Nonce
		canonical, err := Canonicalize(candidateMessage, signature.KeyID, v.algorithm)
		if err != nil {
			return Verification{}, &VerificationError{Kind: ErrInvalidSignature, Diagnostic: "canonical request rejected"}
		}
		expected := signWithHash(v.hash, key.Secret, canonical)
		if hmac.Equal(candidate, expected) {
			return Verification{KeyID: key.ID, Algorithm: v.algorithm, Timestamp: signature.Timestamp, Nonce: signature.Nonce}, nil
		}
	}

	return Verification{}, &VerificationError{Kind: ErrInvalidSignature, Diagnostic: "no supplied signature authenticated"}
}

// VerifyAndRecord authenticates a message and then uses the replay store's
// atomic check-and-record operation. It fails closed on every storage error.
// With no store, eventID is ignored. Otherwise derive it from authenticated
// components: the separate argument is not automatically covered by the MAC.
// Recording occurs before business handling and is not rolled back on its failure.
// Store expiry is sampled verification-clock time plus ReplayTTL, not signed time.
func (v *Verifier) VerifyAndRecord(
	ctx context.Context,
	message Message,
	signatures []Signature,
	eventID string,
) (Verification, error) {
	verification, err := v.Verify(message, signatures)
	if err != nil {
		return Verification{}, err
	}
	if v.replay == nil {
		return verification, nil
	}
	if eventID == "" {
		return Verification{}, &VerificationError{Kind: ErrMissingEventID, Diagnostic: "replay protection requires an event ID"}
	}

	if err := v.recordReplay(ctx, verification, eventID); err != nil {
		return Verification{}, err
	}

	return verification, nil
}

func (v *Verifier) recordReplay(ctx context.Context, verification Verification, eventID string) (err error) {
	started := v.clock()
	defer func() {
		outcome := OutcomeSuccess
		if err != nil {
			outcome = OutcomeRejected
		}
		observeSafely(v.observer, ctx, Observation{
			Operation: OperationReplay,
			Outcome:   outcome,
			Reason:    observationReason(err),
			Duration:  elapsed(v.clock, started),
			Algorithm: verification.Algorithm,
		})
	}()
	if eventID == "" {
		return &VerificationError{Kind: ErrMissingEventID, Diagnostic: "replay protection requires an event ID"}
	}
	key := replayKey(v.namespace, eventID)
	recorded, err := v.replay.CheckAndRecord(ctx, key, v.clock().UTC().Add(v.replayTTL))
	if err != nil {
		return &VerificationError{Kind: ErrReplayStore, Diagnostic: "atomic check-and-record failed"}
	}
	if !recorded {
		return &VerificationError{Kind: ErrReplay, Diagnostic: "event ID was already recorded"}
	}

	return nil
}

// Domain-separated length prefixes avoid delimiter collisions; key ID is
// intentionally excluded so overlapping rotation signatures share one replay
// identity (WEBHOOK-DEC-011, TestReplayIdentitySurvivesSecretRotation).
func replayKey(namespace, eventID string) string {
	digest := sha256.New()
	for _, value := range []string{"go-webhook-replay-v1", namespace, eventID} {
		var size [8]byte
		binary.BigEndian.PutUint64(size[:], uint64(len(value)))
		_, _ = digest.Write(size[:])
		_, _ = digest.Write([]byte(value))
	}

	return hex.EncodeToString(digest.Sum(nil))
}

func activeKey(notBefore, notAfter time.Time, revoked bool, at time.Time) bool {
	return !revoked && (notBefore.IsZero() || !at.Before(notBefore)) && (notAfter.IsZero() || !at.After(notAfter))
}

func hashFactory(algorithm Algorithm) (func() hash.Hash, error) {
	switch algorithm {
	case SHA256:
		return sha256.New, nil
	case SHA512:
		return sha512.New, nil
	default:
		return nil, fmt.Errorf("%w: unsupported algorithm %q", ErrInvalidConfiguration, algorithm)
	}
}

func sign(algorithm Algorithm, secret, message []byte) ([]byte, error) {
	factory, err := hashFactory(algorithm)
	if err != nil {
		return nil, err
	}
	return signWithHash(factory, secret, message), nil
}

func signWithHash(factory func() hash.Hash, secret, message []byte) []byte {
	mac := hmac.New(factory, secret)
	_, _ = mac.Write(message)

	return mac.Sum(nil)
}

func defaultNonce() (string, error) {
	return generateNonce(rand.Reader)
}

func generateNonce(reader io.Reader) (string, error) {
	value := make([]byte, 18)
	if _, err := io.ReadFull(reader, value); err != nil {
		return "", fmt.Errorf("%w: entropy unavailable", ErrNonceGeneration)
	}

	return base64.RawURLEncoding.EncodeToString(value), nil
}

func validNonce(nonce string) bool {
	if nonce == "" {
		return false
	}
	if len(nonce) > maxNonceBytes {
		return false
	}

	return utf8.ValidString(nonce)
}
