package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

// OwnedSchoolID is the school a new grade or subject belongs to.
//
// It is the one the creator administers, resolved from their identity and never
// taken from the request: accepting a school id from the body would let anyone
// file a grade under a school they do not belong to, and grades are what
// courses hang off.
//
// A caller who administers none is refused rather than defaulted. A row with no
// school is invisible to every listing, so creating one would look like it
// worked and produce nothing.
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

// OwnedSchoolIDSelected is OwnedSchoolID for a caller who picked a school.
//
// A superadmin passes on the school they selected, as everywhere else: they
// administer no school by membership, so without this they could create a
// course in an institution but not the grade or subject that course needs. They
// are still refused a school that is not active, and still have to select one —
// there is no sensible default for an operator who belongs to none.
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
