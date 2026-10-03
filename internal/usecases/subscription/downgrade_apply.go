package subscription

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	DowngradeApplyUsecase interface {
		Execute(ctx context.Context, teacherID string, in DowngradeApplyInput) (*DowngradeApplyOutput, apperrors.ApplicationError)
	}

	downgradeApplyUsecase struct {
		contextFactory appcontext.Factory
	}

	DowngradeApplyInput struct {
		Keep []string `json:"keep"`
	}

	DowngradeApplyOutput struct {
		Data DowngradeApplyData `json:"data"`
	}

	DowngradeApplyData struct {
		MaxStudents int      `json:"max_students"`
		Deactivated []string `json:"deactivated"`
	}
)

func NewDowngradeApplyUsecase(contextFactory appcontext.Factory) DowngradeApplyUsecase {
	return &downgradeApplyUsecase{contextFactory: contextFactory}
}

func (u *downgradeApplyUsecase) Execute(ctx context.Context, teacherID string, in DowngradeApplyInput) (*DowngradeApplyOutput, apperrors.ApplicationError) {
	keep := in.Keep
	app := u.contextFactory()
	scope, appErr := scopeFor(ctx, app, "", teacherID)
	if appErr != nil {
		return nil, appErr
	}
	if !scope.Enforced() {
		return &DowngradeApplyOutput{Data: DowngradeApplyData{Deactivated: []string{}}}, nil
	}
	byActivity, err := app.Repositories.School.ListStudentsByActivity(ctx, scope.SchoolID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	deactivated := domain.StudentsToDeactivate(byActivity, scope.Plan.MaxStudents, keep)
	if deactivated == nil {
		deactivated = []string{}
	}
	for _, studentID := range deactivated {
		if err := app.Repositories.School.SetMemberActive(ctx, scope.SchoolID, studentID, false); err != nil {
			return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
		}
	}
	return &DowngradeApplyOutput{Data: DowngradeApplyData{MaxStudents: scope.Plan.MaxStudents, Deactivated: deactivated}}, nil
}
