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

func OwnedSchoolIDSelected(ctx context.Context, app *appcontext.Context, userID, schoolID string) (string, apperrors.ApplicationError) {
	if schoolID == "" {
		return OwnedSchoolID(ctx, app, userID)
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
