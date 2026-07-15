# Replay and idempotency

Authentication proves message integrity; replay protection proves an event ID
has not already been accepted inside a configured namespace and TTL. They are
separate guarantees. Event IDs are extracted only after authentication and
are hashed with namespace and key ID before storage.

`ReplayStore.CheckAndRecord` must atomically return `true` only when it both
observed absence and recorded the expiry. Storage errors reject closed.
Implementations must honor context cancellation and provide tenant-separated
namespaces. The contract deliberately does not promise exactly-once handling.

`adapters/goidempotency` maps this contract to a fail-closed
`go-idempotency.Service` lease. Scope includes namespace, tenant, operation,
caller, and the already-hashed replay value. A lease acquisition is new;
retained, in-progress, replay, and conflict outcomes are replay rejection.
Choose a TTL covering clock tolerance plus the provider's maximum retry age.
