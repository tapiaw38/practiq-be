package identity

import (
	"context"
	"errors"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/authapi"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

// batchSize mirrors auth-api-be's own cap on GET /user/batch.
const batchSize = 200

// Names resolves display identity for a set of usernames in one or more
// round trips to auth-api-be, deduplicated and chunked to the endpoint's
// cap. Unknown ids are simply absent from the result map — callers should
// fall back to the bare id when a lookup misses.
// A rejected token is reported as 401 rather than 500: the status decision
// lives here because this is the only place that sees auth-api-be's own
// response, and every caller used to flatten it into a server error.
func Names(ctx context.Context, client authapi.Client, bearerToken string, ids []string) (map[string]authapi.UserInfo, apperrors.ApplicationError) {
	unique := make(map[string]bool, len(ids))
	deduped := make([]string, 0, len(ids))
	for _, id := range ids {
		if id == "" || unique[id] {
			continue
		}
		unique[id] = true
		deduped = append(deduped, id)
	}

	result := make(map[string]authapi.UserInfo, len(deduped))
	for i := 0; i < len(deduped); i += batchSize {
		end := i + batchSize
		if end > len(deduped) {
			end = len(deduped)
		}
		users, err := client.GetBatch(ctx, bearerToken, deduped[i:end])
		if err != nil {
			var upstream *authapi.UpstreamError
			if errors.As(err, &upstream) && upstream.Unauthorized() {
				return nil, apperrors.NewUnauthorizedError()
			}
			return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
		}
		for _, u := range users {
			result[u.Username] = u
		}
	}
	return result, nil
}

// FullName joins first and last name, falling back to the bare id when the
// lookup missed (unknown id, or auth-api-be call failed upstream).
func FullName(info authapi.UserInfo, fallbackID string) string {
	name := strings.TrimSpace(info.FirstName + " " + info.LastName)
	if name == "" {
		return fallbackID
	}
	return name
}
