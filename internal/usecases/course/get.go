package course

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type GetUsecase interface {
	Execute(context.Context, string) (*CourseOutput, apperrors.ApplicationError)
}

type getUsecase struct {
	factory appcontext.Factory
}

func NewGetUsecase(factory appcontext.Factory) GetUsecase {
	return &getUsecase{factory: factory}
}

func (u *getUsecase) Execute(ctx context.Context, id string) (*CourseOutput, apperrors.ApplicationError) {
	app := u.factory()

	c, err := app.Repositories.Course.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if c == nil {
		return nil, apperrors.NewApplicationError(mappings.CourseNotFoundError, nil)
	}

	return &CourseOutput{Data: toCourseData(*c)}, nil
}
