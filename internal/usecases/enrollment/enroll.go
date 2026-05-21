package enrollment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type EnrollUsecase interface {
	Execute(context.Context, string, string) (*EnrollOutput, apperrors.ApplicationError)
}

type enrollUsecase struct {
	factory appcontext.Factory
}

func NewEnrollUsecase(factory appcontext.Factory) EnrollUsecase {
	return &enrollUsecase{factory: factory}
}

func (u *enrollUsecase) Execute(ctx context.Context, courseID, studentID string) (*EnrollOutput, apperrors.ApplicationError) {
	app := u.factory()

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

	return &EnrollOutput{Message: "enrolled successfully"}, nil
}
