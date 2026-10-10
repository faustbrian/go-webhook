package webhook

import (
	"context"
	"math"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestRetryPolicyDelayRejectsInvalidBounds(t *testing.T) {
	for _, policy := range []RetryPolicy{
		{BaseDelay: -time.Second, MaxDelay: time.Second},
		{BaseDelay: 0, MaxDelay: -time.Second},
		{BaseDelay: 2 * time.Second, MaxDelay: time.Second},
	} {
		for _, retryAfter := range []string{"", "1", "9223372036854775807"} {
			if got := policy.Delay(2, time.Unix(0, 0), retryAfter); got != 0 {
				t.Errorf("invalid policy %#v Retry-After %q: delay = %v, want zero", policy, retryAfter, got)
			}
		}
	}
}

func TestRetryPolicyDelayPreservesFiniteBoundaries(t *testing.T) {
	for _, test := range []struct {
		name       string
		policy     RetryPolicy
		attempt    int
		retryAfter string
		want       time.Duration
	}{
		{
			name:    "equal positive delay bounds",
			policy:  RetryPolicy{BaseDelay: time.Second, MaxDelay: time.Second},
			attempt: 1,
			want:    time.Second,
		},
		{
			name:       "whole seconds below fractional cap",
			policy:     RetryPolicy{BaseDelay: time.Nanosecond, MaxDelay: 1500 * time.Millisecond},
			attempt:    1,
			retryAfter: "1",
			want:       time.Second,
		},
		{
			name:    "largest duration cap preserves small backoff",
			policy:  RetryPolicy{BaseDelay: time.Nanosecond, MaxDelay: time.Duration(math.MaxInt64)},
			attempt: 2,
			want:    2 * time.Nanosecond,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := test.policy.Delay(test.attempt, time.Unix(0, 0), test.retryAfter); got != test.want {
				t.Errorf("Delay() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestRetryPolicyDelayHandlesExtremeAttempts(t *testing.T) {
	const probeEnvironment = "GOLIB_WEBHOOK_DELAY_PROBE"
	if os.Getenv(probeEnvironment) != "1" {
		executable, err := os.Executable()
		if err != nil {
			t.Fatal(err)
		}
		// Isolate nonterminating mutations so they fail this oracle instead of
		// exhausting the whole mutation campaign's test deadline.
		ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
		defer cancel()
		// #nosec G204 -- Own test executable, fixed test selector and explicit bounded child lifetime.
		command := exec.CommandContext(ctx, executable, "-test.run=^TestRetryPolicyDelayHandlesExtremeAttempts$", "-test.count=1")
		command.Env = append(os.Environ(), probeEnvironment+"=1")
		output, err := command.CombinedOutput()
		if ctx.Err() != nil {
			t.Fatalf("Delay extreme-attempt probe exceeded its deadline: %v", ctx.Err())
		}
		if err != nil {
			t.Fatalf("Delay extreme-attempt probe failed: %v\n%s", err, output)
		}
		return
	}
	policy := RetryPolicy{BaseDelay: time.Second, MaxDelay: 5 * time.Second}
	for _, test := range []struct {
		attempt int
		want    time.Duration
	}{
		{math.MinInt, time.Second},
		{0, time.Second},
		{1, time.Second},
		{math.MaxInt, 5 * time.Second},
	} {
		if got := policy.Delay(test.attempt, time.Unix(0, 0), ""); got != test.want {
			t.Errorf("attempt %d: delay = %v, want %v", test.attempt, got, test.want)
		}
	}
	zero := RetryPolicy{BaseDelay: 0, MaxDelay: time.Second}
	if got := zero.Delay(math.MaxInt, time.Unix(0, 0), ""); got != 0 {
		t.Errorf("zero base: delay = %v, want zero", got)
	}
}
