package revocation

import (
	"context"
	"sync"
	"time"
)

type Lookup func(ctx context.Context, bearerToken, userID string) (uint, error)

type entry struct {
	version  uint
	fetchedA time.Time
}

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

func (c *Checker) Forget(userID string) {
	c.mu.Lock()
	delete(c.entries, userID)
	c.mu.Unlock()
}
