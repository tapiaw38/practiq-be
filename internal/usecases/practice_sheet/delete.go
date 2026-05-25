package practicesheet

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type DeleteUsecase interface {
	Execute(context.Context, string) apperrors.ApplicationError
}

type deleteUsecase struct {
	factory appcontext.Factory
}

func NewDeleteUsecase(factory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{factory: factory}
}

func (u *deleteUsecase) Execute(ctx context.Context, id string) apperrors.ApplicationError {
	app := u.factory()

	if err := app.Repositories.PracticeSheet.Delete(ctx, id); err != nil {
		return apperrors.NewApplicationError(mappings.PracticeSheetDeleteError, err)
	}

	return nil
}
