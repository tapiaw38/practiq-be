package practicesheet

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type GetUsecase interface {
	Execute(context.Context, string) (*PracticeSheetOutput, apperrors.ApplicationError)
}

type getUsecase struct {
	factory appcontext.Factory
}

func NewGetUsecase(factory appcontext.Factory) GetUsecase {
	return &getUsecase{factory: factory}
}

func (u *getUsecase) Execute(ctx context.Context, id string) (*PracticeSheetOutput, apperrors.ApplicationError) {
	app := u.factory()

	ps, err := app.Repositories.PracticeSheet.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}
	if ps == nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetNotFoundError, nil)
	}

	return &PracticeSheetOutput{Data: toSheetData(*ps)}, nil
}
