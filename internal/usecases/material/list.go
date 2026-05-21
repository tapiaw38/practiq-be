package material

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ListUsecase interface {
	Execute(context.Context, string) (*MaterialListOutput, apperrors.ApplicationError)
}

type listUsecase struct {
	factory appcontext.Factory
}

func NewListUsecase(factory appcontext.Factory) ListUsecase {
	return &listUsecase{factory: factory}
}

func (u *listUsecase) Execute(ctx context.Context, courseID string) (*MaterialListOutput, apperrors.ApplicationError) {
	app := u.factory()

	materials, err := app.Repositories.Material.List(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.MaterialListError, err)
	}

	var data []MaterialData
	for _, m := range materials {
		data = append(data, toMaterialData(m))
	}
	if data == nil {
		data = []MaterialData{}
	}

	return &MaterialListOutput{Data: data}, nil
}
