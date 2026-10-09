package studentprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	GetStudentProgressUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, studentID string) (*GetStudentProgressOutput, apperrors.ApplicationError)
	}

	getStudentProgressUsecase struct {
		contextFactory appcontext.Factory
	}

	GetStudentProgressOutput struct {
		Data                 []ProgressData `json:"data"`
		LastPracticedSheetID string         `json:"last_practiced_sheet_id,omitempty"`
	}
)

func NewGetStudentProgressUsecase(contextFactory appcontext.Factory) GetStudentProgressUsecase {
	return &getStudentProgressUsecase{contextFactory: contextFactory}
}

func (u *getStudentProgressUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, studentID string) (*GetStudentProgressOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if !isSuperAdmin {
		hasAccess, err := app.Repositories.TeacherStudentAssignment.HasAccess(ctx, requesterID, studentID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.AssignmentListError, err)
		}
		if !hasAccess {
			return nil, apperrors.NewForbiddenError()
		}
	}

	list, err := app.Repositories.StudentProgress.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
	}

	loc := studentLocation(ctx, app, studentID)
	var data []ProgressData
	for _, p := range list {
		data = append(data, toProgressData(p, loc))
	}
	if data == nil {
		data = []ProgressData{}
	}

	return &GetStudentProgressOutput{Data: data}, nil
}
