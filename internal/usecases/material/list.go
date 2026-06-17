package material

import (
	"context"

	materialRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/material"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListUsecase interface {
		Execute(context.Context, ListInput) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListInput struct {
		CourseID string
		Limit    int
		Offset   int
	}

	ListOutput struct {
		Data []MaterialData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, input ListInput) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	filter := materialRepo.ListFilter{
		CourseID: input.CourseID,
		Limit:    input.Limit,
		Offset:   input.Offset,
	}

	materials, err := app.Repositories.Material.List(ctx, filter)
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
