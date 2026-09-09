package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

// EnsureAdministers refuses anyone who is not an admin of the given school.
//
// This is what makes it safe to open grade and subject management beyond the
// platform superadmin: the route says a teacher may ask, and this says which
// rows they may touch. Without it, opening the routes would let any teacher
// rename another school's grades.
//
// A superadmin passes, as everywhere else.
func EnsureAdministers(
	ctx context.Context,
	app *appcontext.Context,
	userID string,
	isSuperAdmin bool,
	schoolID string,
) apperrors.ApplicationError {
	if isSuperAdmin {
		return nil
	}
	// A row with no school predates the split or was created while nothing
	// assigned one. Only a superadmin can touch it, or the first teacher to
	// find it would take it over.
	if schoolID == "" {
		return apperrors.NewForbiddenError()
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
