package grade

import (
	"context"
	"slices"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListGradesByUsersUsecase interface {
		// Execute mirrors ListUserGrades' rule per id: the caller sees their
		// own grades always, and everyone else's only if they teach. Ids the
		// caller can't see are dropped rather than failing the whole batch.
		Execute(ctx context.Context, requesterID string, isTeacher bool, userIDs []string) (*ListGradesByUsersOutput, apperrors.ApplicationError)
	}

	listGradesByUsersUsecase struct {
		contextFactory appcontext.Factory
	}

	ListGradesByUsersOutput struct {
		Data map[string][]GradeData `json:"data"`
	}
)

func NewListGradesByUsersUsecase(contextFactory appcontext.Factory) ListGradesByUsersUsecase {
	return &listGradesByUsersUsecase{contextFactory: contextFactory}
}

// filterAllowedUserIDs applies the same rule as the single-user endpoint to
// every id in the batch: teachers see anyone, everyone else only themselves.
func filterAllowedUserIDs(requesterID string, isTeacher bool, userIDs []string) []string {
	if isTeacher {
		return userIDs
	}
	if slices.Contains(userIDs, requesterID) {
		return []string{requesterID}
	}
	return nil
}

func (u *listGradesByUsersUsecase) Execute(ctx context.Context, requesterID string, isTeacher bool, userIDs []string) (*ListGradesByUsersOutput, apperrors.ApplicationError) {
	allowedIDs := filterAllowedUserIDs(requesterID, isTeacher, userIDs)

	app := u.contextFactory()
	grades, err := app.Repositories.Grade.ListGradesByUsers(ctx, allowedIDs)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeListError, err)
	}

	data := make(map[string][]GradeData, len(grades))
	for userID, userGrades := range grades {
		items := make([]GradeData, 0, len(userGrades))
		for _, grade := range userGrades {
			items = append(items, toGradeData(grade))
		}
		data[userID] = items
	}

	return &ListGradesByUsersOutput{Data: data}, nil
}
