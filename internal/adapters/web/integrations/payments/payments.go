// Package payments talks to the payments service, which owns money: plans,
// subscriptions and whatever the gateway says about them.
//
// It deliberately knows nothing about what a plan allows. The payments service
// stores that as opaque metadata on the plan and hands it back untouched, and
// this package passes it through as a map for the domain to interpret. That is
// what lets Practiq change who counts as a student — a rule that has already
// changed once — without deploying the payments service.
package payments

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type (
	// Plan is one purchasable plan. Limits live in Metadata.
	Plan struct {
		ID            int            `json:"id"`
		Name          string         `json:"name"`
		Description   string         `json:"description"`
		Amount        float64        `json:"amount"`
		Currency      string         `json:"currency"`
		Interval      string         `json:"interval"`
		IntervalCount int            `json:"interval_count"`
		Metadata      map[string]any `json:"metadata"`
		Active        int            `json:"active"`
	}

	// Entitlement answers "is this user paid up, and what does their plan
	// allow" in one call.
	Entitlement struct {
		UserID         string         `json:"user_id"`
		Active         bool           `json:"active"`
		SubscriptionID *int           `json:"subscription_id"`
		PlanID         *int           `json:"plan_id"`
		AccessUntil    *time.Time     `json:"access_until"`
		Metadata       map[string]any `json:"metadata"`
	}

	Client interface {
		// ListPlans returns the plans on offer, cheapest first is not
		// guaranteed — the caller orders them.
		ListPlans(ctx context.Context) ([]Plan, error)
		// GetEntitlement never returns an error for "no subscription": an
		// inactive entitlement is the normal state of a teacher on the free
		// plan, not a failure.
		GetEntitlement(ctx context.Context, userID string) (*Entitlement, error)
	}

	client struct {
		baseURL string
		apiKey  string
		http    *http.Client
	}

	// Unavailable marks a payments outage, so callers can tell "this teacher
	// has no subscription" from "we could not find out".
	Unavailable struct{ Err error }
)

func (e *Unavailable) Error() string { return "payments unavailable: " + e.Err.Error() }
func (e *Unavailable) Unwrap() error { return e.Err }

func NewClient(baseURL, apiKey string) Client {
	return &client{
		baseURL: baseURL,
		apiKey:  apiKey,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *client) ListPlans(ctx context.Context) ([]Plan, error) {
	var plans []Plan
	if err := c.get(ctx, "/api/v1/subscriptions/plans", &plans); err != nil {
		return nil, err
	}
	return plans, nil
}

func (c *client) GetEntitlement(ctx context.Context, userID string) (*Entitlement, error) {
	var entitlement Entitlement
	path := fmt.Sprintf("/api/v1/subscriptions/subscriptions/user/%s/entitlement", url.PathEscape(userID))
	if err := c.get(ctx, path, &entitlement); err != nil {
		return nil, err
	}
	return &entitlement, nil
}

func (c *client) get(ctx context.Context, path string, out any) error {
	if c.baseURL == "" {
		return &Unavailable{Err: fmt.Errorf("payments service url is not configured")}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return &Unavailable{Err: err}
	}
	// The key is what identifies Practiq to the payments service; it also
	// decides which product's rows the answer can contain.
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return &Unavailable{Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &Unavailable{Err: fmt.Errorf("payments responded %d", resp.StatusCode)}
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return &Unavailable{Err: err}
	}
	return nil
}
