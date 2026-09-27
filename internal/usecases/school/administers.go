package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

func EnsureAdministers(
	ctx context.Context,
	app *appcontext.Context,
	userID string,
	isSuperAdmin bool,
	schoolID string,
) apperrors.ApplicationError {

	if schoolID == "" {
		if isSuperAdmin {
			return nil
		}
		return apperrors.NewForbiddenError()
	}

	school, err := app.Repositories.School.Get(ctx, schoolID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	if school == nil || school.Status != domain.SchoolStatusActive {
		return apperrors.NewForbiddenError()
	}
	if isSuperAdmin {
		return nil
	}

	members, err := app.Repositories.School.ListForUser(ctx, userID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	for _, member := range members {
		if member.SchoolID == schoolID && member.Role == domain.SchoolRoleAdmin && member.Active {
			return nil
		}
	}
	return apperrors.NewForbiddenError()
}

func RequireActive(ctx context.Context, app *appcontext.Context, schoolID string) apperrors.ApplicationError {
	if schoolID == "" {
		return apperrors.NewForbiddenError()
	}
	school, err := app.Repositories.School.Get(ctx, schoolID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	if school == nil || school.Status != domain.SchoolStatusActive {
		return apperrors.NewForbiddenError()
	}
	return nil
}
