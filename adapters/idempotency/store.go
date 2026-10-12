// Package webhookidempotency adapts idempotency leases to webhook replay checks.
// New borrows a service and explicit full scope; CheckAndRecord uses fail-closed
// Begin. Backend atomicity/lease expiry/availability remain with idempotency.
// No backend, background worker, completion transaction or close is owned.
package webhookidempotency

import (
	"context"
	"errors"
	"fmt"
	"time"

	idempotency "github.com/faustbrian/go-idempotency"
	webhook "github.com/faustbrian/go-webhook/v3"
)

var (
	// ErrInvalidConfig identifies an absent service or invalid collision scope.
	ErrInvalidConfig = errors.New("webhook/idempotency: invalid configuration")
	// ErrInvalidExpiry means sampled expiry has no positive lease duration.
	ErrInvalidExpiry = errors.New("webhook/idempotency: expiry must be in the future")
)

// Config supplies the durable service and full collision scope.
type Config struct {
	// Service is a mandatory borrowed idempotency service; caller owns backend/lifetime.
	Service *idempotency.Service
	// Namespace is a required validated idempotency scope component.
	Namespace string
	// Tenant is a required validated tenant scope, not inferred from a request.
	Tenant string
	// Operation is a required validated operation scope.
	Operation string
	// Caller is a required validated sender/application scope.
	Caller string
	// Clock is borrowed; nil uses time.Now and expiry is converted to a relative lease.
	Clock func() time.Time
}

// Store implements webhook.ReplayStore through fail-closed leases. New is
// required; zero Store is unusable. Service/clock must support concurrent calls
// and bounded work. No Complete/Release is called; expiry permits later takeover,
// and lease admission does not promise exactly-once business effects.
type Store struct {
	service   *idempotency.Service
	namespace string
	tenant    string
	operation string
	caller    string
	clock     func() time.Time
}

var _ webhook.ReplayStore = (*Store)(nil)

// New validates the durable service and every key scope component.
func New(config Config) (*Store, error) {
	if config.Service == nil {
		return nil, ErrInvalidConfig
	}
	if config.Namespace == "" {
		return nil, ErrInvalidConfig
	}
	if config.Tenant == "" {
		return nil, ErrInvalidConfig
	}
	if config.Operation == "" {
		return nil, ErrInvalidConfig
	}
	if config.Caller == "" {
		return nil, ErrInvalidConfig
	}
	if _, err := idempotency.NewKey(
		config.Namespace,
		config.Tenant,
		config.Operation,
		config.Caller,
		"validation",
	); err != nil {
		return nil, fmt.Errorf("%w: invalid scope: %v", ErrInvalidConfig, err)
	}
	clock := config.Clock
	if clock == nil {
		clock = time.Now
	}

	return &Store{
		service:   config.Service,
		namespace: config.Namespace,
		tenant:    config.Tenant,
		operation: config.Operation,
		caller:    config.Caller,
		clock:     clock,
	}, nil
}

// CheckAndRecord acquires a lease through idempotency. Acquired and stale
// takeover outcomes are new records; every retained matching outcome is a
// replay. Backend unavailability remains an error through fail-closed mode.
func (s *Store) CheckAndRecord(ctx context.Context, value string, expiresAt time.Time) (bool, error) {
	now := s.clock().UTC()
	lease := expiresAt.Sub(now)
	if lease <= 0 {
		return false, ErrInvalidExpiry
	}
	key, err := idempotency.NewKey(s.namespace, s.tenant, s.operation, s.caller, value)
	if err != nil {
		return false, fmt.Errorf("webhook/idempotency: replay key: %w", err)
	}
	// The version is a non-empty compile-time protocol constant, so this
	// constructor cannot fail after dependency contract validation.
	fingerprint, _ := idempotency.NewFingerprint("go-webhook-replay-v1", []byte(value))
	result, err := s.service.Begin(ctx, idempotency.BeginRequest{
		Acquire: idempotency.AcquireRequest{
			Key:         key,
			Fingerprint: fingerprint,
			Lease:       lease,
		},
		Availability: idempotency.AvailabilityFailClosed,
	})
	if err != nil {
		return false, fmt.Errorf("webhook/idempotency: acquire replay lease: %w", err)
	}

	return result.Execute, nil
}
