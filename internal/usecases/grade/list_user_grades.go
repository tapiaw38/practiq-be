package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListUserGradesUsecase interface {
		Execute(context.Context, string) (*ListUserGradesOutput, apperrors.ApplicationError)
	}

	listUserGradesUsecase struct {
		contextFactory appcontext.Factory
	}

	ListUserGradesOutput struct {
		Data []GradeData `json:"data"`
	}
)

func NewListUserGradesUsecase(contextFactory appcontext.Factory) ListUserGradesUsecase {
	return &listUserGradesUsecase{contextFactory: contextFactory}
}

func (u *listUserGradesUsecase) Execute(ctx context.Context, userID string) (*ListUserGradesOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	grades, err := app.Repositories.Grade.ListUserGrades(ctx, userID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeListError, err)
	}
	data := make([]GradeData, 0, len(grades))
	for _, grade := range grades {
		data = append(data, toGradeData(grade))
	}
	return &ListUserGradesOutput{Data: data}, nil
}
