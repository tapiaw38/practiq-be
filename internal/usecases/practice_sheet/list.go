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
		Execute(context.Context, string, bool, ListInput) (*ListOutput, apperrors.ApplicationError)
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

func (u *listUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, input ListInput) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if appErr := requesterCanReadCourse(ctx, app, requesterID, isSuperAdmin, input.CourseID); appErr != nil {
		return nil, appErr
	}

	filter := practiceSheetRepo.ListFilter{
		CourseID: input.CourseID,
		Limit:    input.Limit,
		Offset:   input.Offset,
	}

	sheets, err := app.Repositories.PracticeSheet.List(ctx, filter)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetListError, err)
	}

	// The listing carries no exercise bodies, so there is nothing here that
	// depends on who is asking: the answers and statements it used to gate are
	// no longer part of this payload at all.
	var data []PracticeSheetData
	for _, ps := range sheets {
		data = append(data, toSheetSummary(ps))
	}
	if data == nil {
		data = []PracticeSheetData{}
	}

	return &ListOutput{Data: data}, nil
}
