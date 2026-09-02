package authapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type (
	UserInfo struct {
		Username  string
		FirstName string
		LastName  string
		Email     string
	}

	Client interface {
		// GetByEmail forwards the caller's own bearer token (auth-api-be
		// requires superadmin for this lookup) and returns nil, nil when no
		// account exists for that email — not found is not an error here.
		GetByEmail(ctx context.Context, bearerToken, email string) (*UserInfo, error)
		// GetBatch resolves display identity (name/email) for a set of
		// usernames in one round trip — open to any authenticated caller.
		// Unknown usernames are silently omitted from the result.
		GetBatch(ctx context.Context, bearerToken string, usernames []string) ([]UserInfo, error)
	}

	client struct {
		baseURL string
		http    *http.Client
	}
)

func NewClient(baseURL string) Client {
	return &client{baseURL: baseURL, http: &http.Client{Timeout: 10 * time.Second}}
}

func (c *client) GetByEmail(ctx context.Context, bearerToken, email string) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/user/by-email?email="+url.QueryEscape(email), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", bearerToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode >= 300 {
		var errBody map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return nil, fmt.Errorf("auth-api-be user lookup failed (status %d): %v", resp.StatusCode, errBody)
	}

	var parsed struct {
		Data struct {
			Username  string `json:"username"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Email     string `json:"email"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	return &UserInfo{
		Username:  parsed.Data.Username,
		FirstName: parsed.Data.FirstName,
		LastName:  parsed.Data.LastName,
		Email:     parsed.Data.Email,
	}, nil
}

func (c *client) GetBatch(ctx context.Context, bearerToken string, usernames []string) ([]UserInfo, error) {
	if len(usernames) == 0 {
		return nil, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/user/batch?ids="+url.QueryEscape(strings.Join(usernames, ",")), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", bearerToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		var errBody map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return nil, fmt.Errorf("auth-api-be batch lookup failed (status %d): %v", resp.StatusCode, errBody)
	}

	var parsed struct {
		Data []struct {
			ID        string `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Email     string `json:"email"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	out := make([]UserInfo, 0, len(parsed.Data))
	for _, d := range parsed.Data {
		out = append(out, UserInfo{Username: d.ID, FirstName: d.FirstName, LastName: d.LastName, Email: d.Email})
	}
	return out, nil
}
