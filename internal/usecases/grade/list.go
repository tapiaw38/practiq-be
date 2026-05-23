package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ListUsecase interface {
	Execute(context.Context) (*GradeListOutput, apperrors.ApplicationError)
}

type listUsecase struct {
	factory appcontext.Factory
}

func NewListUsecase(factory appcontext.Factory) ListUsecase {
	return &listUsecase{factory: factory}
}

func (u *listUsecase) Execute(ctx context.Context) (*GradeListOutput, apperrors.ApplicationError) {
	app := u.factory()

	grades, err := app.Repositories.Grade.List(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeListError, err)
	}

	data := make([]GradeData, 0, len(grades))
	for _, grade := range grades {
		data = append(data, toGradeData(grade))
	}

	return &GradeListOutput{Data: data}, nil
}
