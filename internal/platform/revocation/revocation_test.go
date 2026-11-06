package revocation

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestChecker(t *testing.T) {
	ctx := context.Background()

	t.Run("a token matching the live version passes", func(t *testing.T) {
		c := NewChecker(func(context.Context, string, string) (uint, error) { return 3, nil }, time.Minute)
		if !c.Current(ctx, "bearer", "user-1", 3) {
			t.Fatal("expected the token to be accepted")
		}
	})

	t.Run("a token from before a password change is refused", func(t *testing.T) {
		c := NewChecker(func(context.Context, string, string) (uint, error) { return 4, nil }, time.Minute)
		if c.Current(ctx, "bearer", "user-1", 3) {
			t.Fatal("expected the superseded token to be refused")
		}
	})

	t.Run("repeat requests reuse the cached answer", func(t *testing.T) {
		calls := 0
		c := NewChecker(func(context.Context, string, string) (uint, error) {
			calls++
			return 1, nil
		}, time.Minute)

		for i := 0; i < 5; i++ {
			c.Current(ctx, "bearer", "user-1", 1)
		}

		if calls != 1 {
			t.Fatalf("auth-api-be was asked %d times, want 1", calls)
		}
	})

	t.Run("the cached answer is dropped once it is stale", func(t *testing.T) {
		calls := 0
		c := NewChecker(func(context.Context, string, string) (uint, error) {
			calls++
			return 1, nil
		}, time.Minute)

		now := time.Now()
		c.now = func() time.Time { return now }
		c.Current(ctx, "bearer", "user-1", 1)

		now = now.Add(time.Minute + time.Second)
		c.Current(ctx, "bearer", "user-1", 1)

		if calls != 2 {
			t.Fatalf("auth-api-be was asked %d times, want 2", calls)
		}
	})

	t.Run("an unreachable auth-api-be does not lock everybody out", func(t *testing.T) {
		c := NewChecker(func(context.Context, string, string) (uint, error) {
			return 0, errors.New("connection refused")
		}, time.Minute)

		if !c.Current(ctx, "bearer", "user-1", 7) {
			t.Fatal("expected the token to be accepted while the lookup is failing")
		}
	})

	t.Run("a failed lookup is not cached as an answer", func(t *testing.T) {
		calls := 0
		c := NewChecker(func(context.Context, string, string) (uint, error) {
			calls++
			return 0, errors.New("connection refused")
		}, time.Minute)

		c.Current(ctx, "bearer", "user-1", 1)
		c.Current(ctx, "bearer", "user-1", 1)

		if calls != 2 {
			t.Fatalf("a failure was cached: asked %d times, want 2", calls)
		}
	})

	t.Run("each person is tracked separately", func(t *testing.T) {
		c := NewChecker(func(_ context.Context, _, userID string) (uint, error) {
			if userID == "user-1" {
				return 1, nil
			}
			return 9, nil
		}, time.Minute)

		if !c.Current(ctx, "bearer", "user-1", 1) {
			t.Fatal("user-1 should pass")
		}
		if c.Current(ctx, "bearer", "user-2", 1) {
			t.Fatal("user-2 should be refused")
		}
	})
}
