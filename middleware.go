package webhook

import (
	"context"
	"fmt"
	"net/http"
)

// ErrorHook receives internal verification failures. Implementations must not
// log request bodies, signatures, secrets, or unredacted sensitive headers.
// Calls are synchronous and may overlap; bound reporting and classify/redact
// underlying error diagnostics. Middleware contains hook panics, not blocking.
type ErrorHook func(ctx context.Context, err error)

// MiddlewareConfig configures inbound verification middleware.
type MiddlewareConfig struct {
	// Request supplies per-request limits, borrowed metadata and event-ID extraction.
	Request RequestOptions
	// FailureStatus defaults from zero to 401; only 400–599 is accepted.
	FailureStatus int
	// OnError is optional synchronous internal reporting; panics are contained, blocking is not.
	OnError ErrorHook
}

type verificationContextKey struct{}
type verifiedBodyContextKey struct{}

// Middleware authenticates and replay-checks a request before invoking next.
// Failure responses always contain only the stable safe message. Zero status
// defaults to 401; next must be non-nil (typed nil is not detected). Metadata and
// callbacks remain borrowed for handler lifetime. Success adds verification and
// an independent body copy to context; next's effects/panics remain caller-owned.
// Replay admission is not rolled back if next fails. Use a constructed verifier.
func (v *Verifier) Middleware(config MiddlewareConfig, next http.Handler) (http.Handler, error) {
	if next == nil {
		return nil, fmt.Errorf("%w: next handler is required", ErrInvalidConfiguration)
	}
	status := config.FailureStatus
	if status == 0 {
		status = http.StatusUnauthorized
	}
	if status < 400 || status > 599 {
		return nil, fmt.Errorf("%w: failure status must be 4xx or 5xx", ErrInvalidConfiguration)
	}

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		verification, body, err := v.VerifyRequest(request.Context(), request, config.Request)
		if err != nil {
			invokeErrorHook(config.OnError, request.Context(), err)
			http.Error(writer, "webhook verification failed", status)

			return
		}
		ctx := context.WithValue(request.Context(), verificationContextKey{}, verification)
		ctx = context.WithValue(ctx, verifiedBodyContextKey{}, append([]byte(nil), body...))
		next.ServeHTTP(writer, request.WithContext(ctx))
	}), nil
}

func invokeErrorHook(hook ErrorHook, ctx context.Context, err error) {
	defer func() { _ = recover() }()
	if hook != nil {
		hook(ctx, err)
	}
}

// VerificationFromContext returns the middleware verification value, or zero/false
// if absent. ctx must be non-nil; the value is attribution, not authorization.
func VerificationFromContext(ctx context.Context) (Verification, bool) {
	verification, ok := ctx.Value(verificationContextKey{}).(Verification)

	return verification, ok
}

// VerifiedBodyFromContext returns caller-owned authenticated body bytes or
// nil/false when absent. ctx must be non-nil; sensitive bytes are not erased.
func VerifiedBodyFromContext(ctx context.Context) ([]byte, bool) {
	body, ok := ctx.Value(verifiedBodyContextKey{}).([]byte)
	if !ok {
		return nil, false
	}

	return append([]byte(nil), body...), true
}
