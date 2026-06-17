package practicesheet

import (
	"context"

	practiceSheetRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/practice_sheet"
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
		Data []PracticeSheetData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, input ListInput) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	filter := practiceSheetRepo.ListFilter{
		CourseID: input.CourseID,
		Limit:    input.Limit,
		Offset:   input.Offset,
	}

	sheets, err := app.Repositories.PracticeSheet.List(ctx, filter)
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
