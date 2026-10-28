package identity

import (
	"context"
	"errors"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/authapi"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

const batchSize = 200

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

func ByEmail(ctx context.Context, client authapi.Client, bearerToken, email string) (*authapi.UserInfo, apperrors.ApplicationError) {
	info, err := client.GetByEmail(ctx, bearerToken, email)
	if err != nil {
		var upstream *authapi.UpstreamError
		if errors.As(err, &upstream) && upstream.Unauthorized() {
			return nil, apperrors.NewUnauthorizedError()
		}
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	return info, nil
}

func FullName(info authapi.UserInfo, fallbackID string) string {
	name := strings.TrimSpace(info.FirstName + " " + info.LastName)
	if name == "" {
		return fallbackID
	}
	return name
}
