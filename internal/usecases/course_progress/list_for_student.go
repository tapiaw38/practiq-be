package courseprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListForStudentUsecase interface {
		Execute(ctx context.Context, requesterID, studentID string, isAdmin bool) (*ListForStudentOutput, apperrors.ApplicationError)
	}

	listForStudentUsecase struct {
		contextFactory appcontext.Factory
	}

	ListForStudentOutput struct {
		Data []CourseProgressData `json:"data"`
	}
)

func NewListForStudentUsecase(contextFactory appcontext.Factory) ListForStudentUsecase {
	return &listForStudentUsecase{contextFactory: contextFactory}
}

func (u *listForStudentUsecase) Execute(ctx context.Context, requesterID, studentID string, isAdmin bool) (*ListForStudentOutput, apperrors.ApplicationError) {
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

	progressList, err := app.Repositories.CourseProgress.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseProgressListError, err)
	}

	data := make([]CourseProgressData, 0, len(progressList))
	for _, p := range progressList {
		data = append(data, toProgressData(p))
	}

	return &ListForStudentOutput{Data: data}, nil
}
