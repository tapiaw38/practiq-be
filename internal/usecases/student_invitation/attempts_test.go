package studentinvitation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAttemptLimiter(t *testing.T) {
	now := time.Now()

	t.Run("allows up to the cap and then blocks", func(t *testing.T) {
		l := &attemptLimiter{attempts: make(map[string][]time.Time)}

		for i := range maxAttemptsPerWindow {
			assert.True(t, l.allow("student-1", now), "attempt %d should be allowed", i+1)
			l.recordFailure("student-1", now)
		}

		assert.False(t, l.allow("student-1", now))
	})

	t.Run("forgets attempts older than the window", func(t *testing.T) {
		l := &attemptLimiter{attempts: make(map[string][]time.Time)}

		old := now.Add(-attemptWindow - time.Minute)
		for range maxAttemptsPerWindow {
			l.recordFailure("student-1", old)
		}

		assert.True(t, l.allow("student-1", now))
		assert.NotContains(t, l.attempts, "student-1",
			"an entry with no recent attempts should be dropped, not kept empty")
	})

	t.Run("counts each student separately", func(t *testing.T) {
		l := &attemptLimiter{attempts: make(map[string][]time.Time)}

		for range maxAttemptsPerWindow {
			l.recordFailure("student-1", now)
		}

		assert.False(t, l.allow("student-1", now))
		assert.True(t, l.allow("student-2", now))
	})

	t.Run("clears the history of whoever got it right", func(t *testing.T) {
		l := &attemptLimiter{attempts: make(map[string][]time.Time)}

		for range maxAttemptsPerWindow {
			l.recordFailure("student-1", now)
		}
		l.clear("student-1")

		assert.True(t, l.allow("student-1", now))
		assert.Empty(t, l.attempts["student-1"])
	})
}
