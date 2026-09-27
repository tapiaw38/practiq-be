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
