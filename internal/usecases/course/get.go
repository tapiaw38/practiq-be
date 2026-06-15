package course

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	GetUsecase interface {
		Execute(context.Context, string) (*GetOutput, apperrors.ApplicationError)
	}

	getUsecase struct {
		contextFactory appcontext.Factory
	}

	GetOutput struct {
		Data CourseData `json:"data"`
	}
)

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{contextFactory: contextFactory}
}

func (u *getUsecase) Execute(ctx context.Context, id string) (*GetOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	c, err := app.Repositories.Course.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if c == nil {
		return nil, apperrors.NewApplicationError(mappings.CourseNotFoundError, nil)
	}

	return &GetOutput{Data: toCourseData(*c)}, nil
}
