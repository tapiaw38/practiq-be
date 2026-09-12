package identity

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/authapi"
)

type stubClient struct {
	users []authapi.UserInfo
	err   error
}

func (s stubClient) GetByEmail(context.Context, string, string) (*authapi.UserInfo, error) {
	return nil, nil
}

func (s stubClient) GetBatch(context.Context, string, []string) ([]authapi.UserInfo, error) {
	return s.users, s.err
}

func TestNamesStatusMapping(t *testing.T) {
	cases := map[string]struct {
		err      error
		expected int
	}{
		"a rejected token is the caller's problem": {
			err:      &authapi.UpstreamError{Status: http.StatusUnauthorized, Detail: "token version mismatch"},
			expected: http.StatusUnauthorized,
		},
		"a forbidden lookup is too": {
			err:      &authapi.UpstreamError{Status: http.StatusForbidden},
			expected: http.StatusUnauthorized,
		},
		"auth-api-be failing is ours": {
			err:      &authapi.UpstreamError{Status: http.StatusInternalServerError},
			expected: http.StatusInternalServerError,
		},
		"so is not reaching it at all": {
			err:      errors.New("connection refused"),
			expected: http.StatusInternalServerError,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, appErr := Names(context.Background(), stubClient{err: tc.err}, "token", []string{"someone"})
			if appErr == nil {
				t.Fatal("expected an error, got none")
			}
			if appErr.StatusCode() != tc.expected {
				t.Fatalf("status = %d, want %d", appErr.StatusCode(), tc.expected)
			}
		})
	}
}

func TestNamesResolvesByUsername(t *testing.T) {
	client := stubClient{users: []authapi.UserInfo{{Username: "nymia", FirstName: "Nymia", LastName: "Tapia", Email: "n@example.com"}}}

	names, appErr := Names(context.Background(), client, "token", []string{"nymia", "nymia", ""})
	if appErr != nil {
		t.Fatalf("unexpected error: %v", appErr)
	}
	if got := FullName(names["nymia"], "nymia"); got != "Nymia Tapia" {
		t.Fatalf("name = %q, want %q", got, "Nymia Tapia")
	}
	if got := FullName(names["ausente"], "ausente"); got != "ausente" {
		t.Fatalf("missing lookup should fall back to the id, got %q", got)
	}
}
