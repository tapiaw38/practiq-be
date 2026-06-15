package material

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListUsecase interface {
		Execute(context.Context, string) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data []MaterialData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, courseID string) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

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

	return &ListOutput{Data: data}, nil
}
