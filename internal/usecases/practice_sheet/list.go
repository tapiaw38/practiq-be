package practicesheet

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
		Data []PracticeSheetData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, courseID string) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

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

	return &ListOutput{Data: data}, nil
}
