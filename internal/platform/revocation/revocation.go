package revocation

import (
	"context"
	"sync"
	"time"
)

// Lookup reports the token version auth-api-be currently holds for a user.
type Lookup func(ctx context.Context, bearerToken, userID string) (uint, error)

type entry struct {
	version  uint
	fetchedA time.Time
}

// Checker answers "is this token still current?" without asking auth-api-be
// on every request.
//
// A password change bumps the user's token version there, which is what
// makes an access token stop working before it expires. Asking upstream on
// every request would put auth-api-be in the path of everything; caching
// for a short while bounds how long a revoked token keeps working to the
// TTL rather than to the token's full lifetime.
type Checker struct {
	lookup Lookup
	ttl    time.Duration

	mu      sync.RWMutex
	entries map[string]entry
	now     func() time.Time
}

func NewChecker(lookup Lookup, ttl time.Duration) *Checker {
	return &Checker{
		lookup:  lookup,
		ttl:     ttl,
		entries: make(map[string]entry),
		now:     time.Now,
	}
}

// Current reports whether the presented version is the live one. When
// auth-api-be cannot be reached the token is accepted: its signature and
// expiry were already checked, and refusing everything would turn one
// service being unreachable into the whole platform being down.
func (c *Checker) Current(ctx context.Context, bearerToken, userID string, presented uint) bool {
	if cached, ok := c.cached(userID); ok {
		return cached == presented
	}

	version, err := c.lookup(ctx, bearerToken, userID)
	if err != nil {
		return true
	}

	c.mu.Lock()
	c.entries[userID] = entry{version: version, fetchedA: c.now()}
	c.mu.Unlock()

	return version == presented
}

func (c *Checker) cached(userID string) (uint, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	found, ok := c.entries[userID]
	if !ok || c.now().Sub(found.fetchedA) >= c.ttl {
		return 0, false
	}
	return found.version, true
}

// Forget drops a cached answer so the next request re-reads it.
func (c *Checker) Forget(userID string) {
	c.mu.Lock()
	delete(c.entries, userID)
	c.mu.Unlock()
}
