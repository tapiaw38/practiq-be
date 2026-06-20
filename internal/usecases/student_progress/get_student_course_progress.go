package studentprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	GetStudentCourseProgressUsecase interface {
		Execute(ctx context.Context, requesterID string, isAdmin bool, studentID, courseID string) (*GetStudentCourseProgressOutput, apperrors.ApplicationError)
	}

	getStudentCourseProgressUsecase struct {
		contextFactory appcontext.Factory
	}

	GetStudentCourseProgressOutput struct {
		Data                 []ProgressData `json:"data"`
		LastPracticedSheetID string         `json:"last_practiced_sheet_id,omitempty"`
	}
)

func NewGetStudentCourseProgressUsecase(contextFactory appcontext.Factory) GetStudentCourseProgressUsecase {
	return &getStudentCourseProgressUsecase{contextFactory: contextFactory}
}

func (u *getStudentCourseProgressUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, studentID, courseID string) (*GetStudentCourseProgressOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if !isAdmin {
		hasAccess, err := app.Repositories.TeacherStudentAssignment.HasAccess(ctx, requesterID, studentID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.AssignmentListError, err)
		}
		if !hasAccess {
			return nil, apperrors.NewForbiddenError()
		}
	}

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

	return &GetStudentCourseProgressOutput{Data: data}, nil
}
