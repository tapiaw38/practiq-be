package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

func OwnedSchoolID(ctx context.Context, app *appcontext.Context, userID string) (string, apperrors.ApplicationError) {
	members, err := app.Repositories.School.ListForUser(ctx, userID)
	if err != nil {
		return "", apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	for _, member := range members {
		if member.Role == domain.SchoolRoleAdmin && member.Active {
			return member.SchoolID, nil
		}
	}
	return "", apperrors.NewBadRequestError("you do not administer a school")
}

func OwnedSchoolIDSelected(ctx context.Context, app *appcontext.Context, userID string, isSuperAdmin bool, schoolID string) (string, apperrors.ApplicationError) {
	if schoolID == "" {
		if isSuperAdmin {
			return "", apperrors.NewBadRequestError("select a school first")
		}
		return OwnedSchoolID(ctx, app, userID)
	}
	if isSuperAdmin {
		if appErr := RequireActive(ctx, app, schoolID); appErr != nil {
			return "", appErr
		}
		return schoolID, nil
	}
	members, err := app.Repositories.School.ListForUser(ctx, userID)
	if err != nil {
		return "", apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	for _, member := range members {
		if member.SchoolID == schoolID && member.Role == domain.SchoolRoleAdmin && member.Active {
			return schoolID, nil
		}
	}
	return "", apperrors.NewBadRequestError("you do not administer this school")
}
