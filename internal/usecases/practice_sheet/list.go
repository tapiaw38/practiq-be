package practicesheet

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ListUsecase interface {
	Execute(context.Context, string) (*PracticeSheetListOutput, apperrors.ApplicationError)
}

type listUsecase struct {
	factory appcontext.Factory
}

func NewListUsecase(factory appcontext.Factory) ListUsecase {
	return &listUsecase{factory: factory}
}

func (u *listUsecase) Execute(ctx context.Context, courseID string) (*PracticeSheetListOutput, apperrors.ApplicationError) {
	app := u.factory()

	sheets, err := app.Repositories.PracticeSheet.List(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetListError, err)
	}

	var data []PracticeSheetData
	for _, ps := range sheets {
		data = append(data, toSheetData(ps))
	}
	if data == nil {
		data = []PracticeSheetData{}
	}

	return &PracticeSheetListOutput{Data: data}, nil
}
