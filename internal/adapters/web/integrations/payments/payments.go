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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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
		AccessUntil    *Timestamp     `json:"access_until"`
		Metadata       map[string]any `json:"metadata"`
	}

	// PlanInput is what a superadmin can set. Amount and limits are separate
	// concerns: the first is what the gateway charges, the second is what
	// Practiq reads out of Metadata.
	PlanInput struct {
		Name        string         `json:"name,omitempty"`
		Description string         `json:"description,omitempty"`
		Amount      *float64       `json:"amount,omitempty"`
		Currency    string         `json:"currency,omitempty"`
		Interval    string         `json:"interval,omitempty"`
		Metadata    map[string]any `json:"metadata,omitempty"`
		Active      *bool          `json:"active,omitempty"`
	}

	SubscriptionInput struct {
		PlanID      int    `json:"plan_id"`
		UserID      string `json:"user_id"`
		PayerEmail  string `json:"payer_email"`
		CardTokenID string `json:"card_token_id"`
	}

	Subscription struct {
		ID     int    `json:"id"`
		PlanID int    `json:"plan_id"`
		UserID string `json:"user_id"`
		Status string `json:"status"`
	}

	// HostedSubscriptionInput subscribes somebody who is not giving us a card,
	// because they intend to pay with their Mercado Pago balance.
	HostedSubscriptionInput struct {
		PlanID     int    `json:"plan_id"`
		UserID     string `json:"user_id"`
		PayerEmail string `json:"payer_email"`
	}

	// HostedSubscription is an agreement waiting for the payer to authorise it
	// at the gateway. Nothing is charged until they do.
	HostedSubscription struct {
		SubscriptionID int    `json:"subscription_id"`
		Status         string `json:"status"`
		InitPoint      string `json:"init_point"`
	}

	Client interface {
		// ListPlans returns the plans on offer, cheapest first is not
		// guaranteed — the caller orders them.
		ListPlans(ctx context.Context) ([]Plan, error)
		// GetEntitlement never returns an error for "no subscription": an
		// inactive entitlement is the normal state of a teacher on the free
		// plan, not a failure.
		GetEntitlement(ctx context.Context, userID string) (*Entitlement, error)
		// ListSubscriptions returns every subscription a user has, whatever
		// its status. GetEntitlement only reports live ones, so a paused
		// subscription is invisible to it and could never be resumed.
		ListSubscriptions(ctx context.Context, userID string) ([]Subscription, error)
		// CreateSubscription hands the gateway a card token the browser
		// produced. The card itself never reaches us.
		CreateSubscription(ctx context.Context, in SubscriptionInput) (*Subscription, error)
		// StartHostedSubscription returns somewhere to send the payer so they
		// can authorise the agreement at the gateway, where their account
		// balance is an option a card form cannot offer.
		StartHostedSubscription(ctx context.Context, in HostedSubscriptionInput) (*HostedSubscription, error)
		CreatePlan(ctx context.Context, in PlanInput) (*Plan, error)
		UpdatePlan(ctx context.Context, planID int, in PlanInput) (*Plan, error)
		// DeactivatePlan takes a plan off the shelf. Subscriptions to it keep
		// working, which is why nothing here deletes a plan.
		DeactivatePlan(ctx context.Context, planID int) (*Plan, error)
		// PauseSubscription stops the charges without ending the agreement, so
		// it can be undone. CancelSubscription cannot.
		PauseSubscription(ctx context.Context, subscriptionID int) (*Subscription, error)
		ResumeSubscription(ctx context.Context, subscriptionID int) (*Subscription, error)
		CancelSubscription(ctx context.Context, subscriptionID int) error
	}

	client struct {
		baseURL string
		apiKey  string
		http    *http.Client
	}

	// Timestamp reads an instant whether or not the payments service says which
	// zone it is in.
	//
	// It serialises naive datetimes, so access_until arrives as
	// "2026-10-24T02:49:15". time.Time insists on an offset and fails the whole
	// decode, which turned a paid subscription into a payments outage and
	// showed the teacher the free plan. A timestamp is not worth that.
	Timestamp struct{ time.Time }

	// Unavailable marks a payments outage, so callers can tell "this teacher
	// has no subscription" from "we could not find out".
	Unavailable struct{ Err error }

	// Rejected is an expected 4xx answer from the payments product. It carries
	// its public error code so the use case can guide a teacher without
	// pretending that Mercado Pago is down.
	Rejected struct {
		StatusCode int
		Code       string
		Message    string
	}
)

// timestampLayouts are tried in order. The offset-bearing one first, so a
// service that does say the zone is believed rather than reinterpreted.
var timestampLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999",
	"2006-01-02T15:04:05",
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	raw := strings.Trim(string(data), `"`)
	if raw == "" || raw == "null" {
		return nil
	}
	for _, layout := range timestampLayouts {
		// Naive datetimes are UTC: that is what the payments service stores,
		// and reading them as local time would move an expiry by hours.
		if parsed, err := time.ParseInLocation(layout, raw, time.UTC); err == nil {
			t.Time = parsed
			return nil
		}
	}
	return fmt.Errorf("payments: unrecognised timestamp %q", raw)
}

func (e *Unavailable) Error() string { return "payments unavailable: " + e.Err.Error() }
func (e *Unavailable) Unwrap() error { return e.Err }
func (e *Rejected) Error() string {
	return fmt.Sprintf("payments rejected %d: %s", e.StatusCode, e.Code)
}

func IsRejected(err error) (*Rejected, bool) {
	var rejected *Rejected
	return rejected, errors.As(err, &rejected)
}

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

func (c *client) ListSubscriptions(ctx context.Context, userID string) ([]Subscription, error) {
	var subscriptions []Subscription
	path := fmt.Sprintf("/api/v1/subscriptions/subscriptions/user/%s", url.PathEscape(userID))
	if err := c.get(ctx, path, &subscriptions); err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (c *client) CreateSubscription(ctx context.Context, in SubscriptionInput) (*Subscription, error) {
	var subscription Subscription
	if err := c.send(ctx, http.MethodPost, "/api/v1/subscriptions/subscriptions", in, &subscription); err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (c *client) StartHostedSubscription(ctx context.Context, in HostedSubscriptionInput) (*HostedSubscription, error) {
	var hosted HostedSubscription
	if err := c.send(ctx, http.MethodPost, "/api/v1/subscriptions/subscriptions/hosted", in, &hosted); err != nil {
		return nil, err
	}
	return &hosted, nil
}

func (c *client) CreatePlan(ctx context.Context, in PlanInput) (*Plan, error) {
	var plan Plan
	if err := c.send(ctx, http.MethodPost, "/api/v1/subscriptions/plans", in, &plan); err != nil {
		return nil, err
	}
	return &plan, nil
}

func (c *client) UpdatePlan(ctx context.Context, planID int, in PlanInput) (*Plan, error) {
	var plan Plan
	path := fmt.Sprintf("/api/v1/subscriptions/plans/%d", planID)
	if err := c.send(ctx, http.MethodPut, path, in, &plan); err != nil {
		return nil, err
	}
	return &plan, nil
}

func (c *client) DeactivatePlan(ctx context.Context, planID int) (*Plan, error) {
	var plan Plan
	path := fmt.Sprintf("/api/v1/subscriptions/plans/%d", planID)
	if err := c.send(ctx, http.MethodDelete, path, nil, &plan); err != nil {
		return nil, err
	}
	return &plan, nil
}

func (c *client) PauseSubscription(ctx context.Context, subscriptionID int) (*Subscription, error) {
	return c.subscriptionAction(ctx, subscriptionID, "pause")
}

func (c *client) ResumeSubscription(ctx context.Context, subscriptionID int) (*Subscription, error) {
	return c.subscriptionAction(ctx, subscriptionID, "resume")
}

func (c *client) CancelSubscription(ctx context.Context, subscriptionID int) error {
	path := fmt.Sprintf("/api/v1/subscriptions/subscriptions/%d/cancel", subscriptionID)
	return c.send(ctx, http.MethodPost, path, nil, nil)
}

func (c *client) subscriptionAction(ctx context.Context, subscriptionID int, action string) (*Subscription, error) {
	var subscription Subscription
	path := fmt.Sprintf("/api/v1/subscriptions/subscriptions/%d/%s", subscriptionID, action)
	if err := c.send(ctx, http.MethodPost, path, nil, &subscription); err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (c *client) send(ctx context.Context, method, path string, body any, out any) error {
	if c.baseURL == "" {
		return &Unavailable{Err: fmt.Errorf("payments service url is not configured")}
	}

	var payload io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return &Unavailable{Err: err}
		}
		payload = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, payload)
	if err != nil {
		return &Unavailable{Err: err}
	}
	req.Header.Set("X-API-Key", c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return &Unavailable{Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return c.responseError(resp)
	}
	if out == nil {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return &Unavailable{Err: err}
	}
	return nil
}

func (c *client) responseError(resp *http.Response) error {
	// Payments only returns our own structured code/message. Read a bounded
	// body anyway: gateway error pages must never become a memory risk.
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
	if resp.StatusCode < 400 || resp.StatusCode >= 500 {
		return &Unavailable{Err: fmt.Errorf("payments responded %d", resp.StatusCode)}
	}
	var envelope struct {
		Detail struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"detail"`
	}
	_ = json.Unmarshal(body, &envelope)
	return &Rejected{
		StatusCode: resp.StatusCode,
		Code:       strings.TrimSpace(envelope.Detail.Code),
		Message:    strings.TrimSpace(envelope.Detail.Message),
	}
}
