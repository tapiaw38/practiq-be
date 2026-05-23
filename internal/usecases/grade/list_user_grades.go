package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ListUserGradesUsecase interface {
	Execute(context.Context, string) (*GradeListOutput, apperrors.ApplicationError)
}

type listUserGradesUsecase struct {
	factory appcontext.Factory
}

func NewListUserGradesUsecase(factory appcontext.Factory) ListUserGradesUsecase {
	return &listUserGradesUsecase{factory: factory}
}

func (u *listUserGradesUsecase) Execute(ctx context.Context, userID string) (*GradeListOutput, apperrors.ApplicationError) {
	app := u.factory()
	grades, err := app.Repositories.Grade.ListUserGrades(ctx, userID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeListError, err)
	}
	data := make([]GradeData, 0, len(grades))
	for _, grade := range grades {
		data = append(data, toGradeData(grade))
	}
	return &GradeListOutput{Data: data}, nil
}
