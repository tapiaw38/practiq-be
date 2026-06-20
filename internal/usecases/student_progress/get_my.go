package studentprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	GetMyProgressUsecase interface {
		Execute(context.Context, string) (*GetMyProgressOutput, apperrors.ApplicationError)
	}

	getMyProgressUsecase struct {
		contextFactory appcontext.Factory
	}

	GetMyProgressOutput struct {
		Data                 []ProgressData `json:"data"`
		LastPracticedSheetID string         `json:"last_practiced_sheet_id,omitempty"`
	}
)

func NewGetMyProgressUsecase(contextFactory appcontext.Factory) GetMyProgressUsecase {
	return &getMyProgressUsecase{contextFactory: contextFactory}
}

func (u *getMyProgressUsecase) Execute(ctx context.Context, studentID string) (*GetMyProgressOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	list, err := app.Repositories.StudentProgress.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
	}

	var data []ProgressData
	for _, p := range list {
		data = append(data, toProgressData(p))
	}
	if data == nil {
		data = []ProgressData{}
	}

	lastSheetID, err := app.Repositories.StudentAttempt.GetLastPracticedSheetID(ctx, studentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AttemptGetError, err)
	}

	return &GetMyProgressOutput{
		Data:                 data,
		LastPracticedSheetID: lastSheetID,
	}, nil
}
