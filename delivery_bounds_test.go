package webhook

import (
	"math"
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

func TestRetryPolicyDelayHandlesExtremeAttempts(t *testing.T) {
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
