package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListUsecase interface {
		Execute(context.Context) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data []GradeData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	grades, err := app.Repositories.Grade.List(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeListError, err)
	}

	data := make([]GradeData, 0, len(grades))
	for _, grade := range grades {
		data = append(data, toGradeData(grade))
	}

	return &ListOutput{Data: data}, nil
}
