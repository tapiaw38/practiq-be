package studentprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	GetCourseProgressUsecase interface {
		Execute(context.Context, string, string) (*GetCourseProgressOutput, apperrors.ApplicationError)
	}

	getCourseProgressUsecase struct {
		contextFactory appcontext.Factory
	}

	GetCourseProgressOutput struct {
		Data                 []ProgressData `json:"data"`
		LastPracticedSheetID string         `json:"last_practiced_sheet_id,omitempty"`
	}
)

func NewGetCourseProgressUsecase(contextFactory appcontext.Factory) GetCourseProgressUsecase {
	return &getCourseProgressUsecase{contextFactory: contextFactory}
}

func (u *getCourseProgressUsecase) Execute(ctx context.Context, studentID, courseID string) (*GetCourseProgressOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	list, err := app.Repositories.StudentProgress.ListByStudentAndCourse(ctx, studentID, courseID)
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

	return &GetCourseProgressOutput{Data: data}, nil
}
