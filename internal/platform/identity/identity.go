package identity

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/authapi"
)

// batchSize mirrors auth-api-be's own cap on GET /user/batch.
const batchSize = 200

// Names resolves display identity for a set of usernames in one or more
// round trips to auth-api-be, deduplicated and chunked to the endpoint's
// cap. Unknown ids are simply absent from the result map — callers should
// fall back to the bare id when a lookup misses.
func Names(ctx context.Context, client authapi.Client, bearerToken string, ids []string) (map[string]authapi.UserInfo, error) {
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
			return nil, err
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
