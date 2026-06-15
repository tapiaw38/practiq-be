package enrollment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	EnrollUsecase interface {
		Execute(context.Context, string, string) (*EnrollOutput, apperrors.ApplicationError)
	}

	enrollUsecase struct {
		contextFactory appcontext.Factory
	}

	EnrollOutput struct {
		Data OperationResultData `json:"data"`
	}
)

func NewEnrollUsecase(contextFactory appcontext.Factory) EnrollUsecase {
	return &enrollUsecase{contextFactory: contextFactory}
}

func (u *enrollUsecase) Execute(ctx context.Context, courseID, studentID string) (*EnrollOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	exists, err := app.Repositories.Enrollment.Exists(ctx, courseID, studentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentCreateError, err)
	}
	if exists {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentAlreadyExistsError, nil)
	}

	if err := app.Repositories.Enrollment.Create(ctx, domain.Enrollment{
		CourseID:  courseID,
		StudentID: studentID,
		Status:    "active",
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentCreateError, err)
	}

	return &EnrollOutput{Data: toOperationResultData(domain.OperationResult{Message: "enrolled successfully"})}, nil
}
