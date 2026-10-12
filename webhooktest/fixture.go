// Package webhooktest provides deterministic clocks, nonces, identifiers,
// signers, and verifiers for consumer tests. Start with NewFixture; predictable
// nonces/IDs are not production randomness or replay protection. No network,
// background worker or process-global clock is owned.
package webhooktest

import (
	"errors"
	"fmt"
	"sync"
	"time"

	webhook "github.com/faustbrian/go-webhook/v3"
)

// ErrInvalidConfig identifies invalid deterministic fixture configuration.
var ErrInvalidConfig = errors.New("webhooktest: invalid fixture configuration")

// Clock is a concurrency-safe manually controlled clock. The zero value returns
// zero time; Set/Advance/Now coordinate one lock and may run concurrently. No
// monotonic clock guarantee, timer or automatic time advancement is provided.
type Clock struct {
	mu  sync.RWMutex
	now time.Time
}

// NewClock constructs a clock at the supplied instant.
func NewClock(now time.Time) *Clock {
	return &Clock{now: now.UTC()}
}

// Now returns the controlled instant.
func (c *Clock) Now() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.now
}

// Set replaces the controlled instant.
func (c *Clock) Set(now time.Time) {
	c.mu.Lock()
	c.now = now.UTC()
	c.mu.Unlock()
}

// Advance moves the controlled instant by duration.
func (c *Clock) Advance(duration time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(duration)
	c.mu.Unlock()
}

// FixtureConfig defines one deterministic signing and verification pair.
type FixtureConfig struct {
	// Algorithm must be a core supported HMAC algorithm.
	Algorithm webhook.Algorithm
	// KeyID must be nonempty for the constructed core pair.
	KeyID string
	// Secret must be nonempty; each core constructor copies it. Never use live credentials in fixtures.
	Secret []byte
	// Time must be nonzero and initializes the shared UTC clock.
	Time time.Time
	// Tolerance must be nonnegative core verifier skew.
	Tolerance time.Duration
	// NoncePrefix must be nonempty; generated values are validated only when signing.
	NoncePrefix string
	// IDPrefix must be nonempty; it prefixes a separate deterministic uint64 sequence.
	IDPrefix string
}

// Fixture owns a shared clock, signer, verifier, and deterministic ID source.
// Construct with NewFixture; zero is unusable. Sequence order follows lock
// acquisition rather than goroutine launch, and uint64 counters eventually wrap.
type Fixture struct {
	// Clock is the shared concurrency-safe manually controlled clock.
	Clock *Clock
	// Signer is the constructed signer using predictable fixture nonces, not production entropy.
	Signer *webhook.Signer
	// Verifier is the constructed verifier with replay protection disabled.
	Verifier *webhook.Verifier
	ids      *sequence
}

// NewFixture constructs a deterministic pair without global state.
func NewFixture(config FixtureConfig) (*Fixture, error) {
	if config.Time.IsZero() || config.NoncePrefix == "" || config.IDPrefix == "" {
		return nil, ErrInvalidConfig
	}
	clock := NewClock(config.Time)
	nonces := &sequence{prefix: config.NoncePrefix}
	signer, err := webhook.NewSigner(webhook.SignerConfig{
		Algorithm:      config.Algorithm,
		Keys:           []webhook.SigningKey{{ID: config.KeyID, Secret: config.Secret}},
		Clock:          clock.Now,
		NonceGenerator: nonces.next,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: signer: %v", ErrInvalidConfig, err)
	}
	verifier, err := webhook.NewVerifier(webhook.VerifierConfig{
		Algorithm: config.Algorithm,
		Keys:      []webhook.VerificationKey{{ID: config.KeyID, Secret: config.Secret}},
		Clock:     clock.Now,
		Tolerance: config.Tolerance,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: verifier: %v", ErrInvalidConfig, err)
	}

	return &Fixture{Clock: clock, Signer: signer, Verifier: verifier, ids: &sequence{prefix: config.IDPrefix}}, nil
}

// IDGenerator returns the next deterministic identifier and matches the
// DeliveryConfig IDGenerator signature.
func (f *Fixture) IDGenerator() (string, error) {
	return f.ids.next()
}

type sequence struct {
	mu     sync.Mutex
	prefix string
	nextID uint64
}

func (s *sequence) next() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++

	return fmt.Sprintf("%s-%d", s.prefix, s.nextID), nil
}
