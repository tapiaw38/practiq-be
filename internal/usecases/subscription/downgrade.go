package subscription

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	// DowngradeUsecase brings a school back inside its plan.
	//
	// It is applied when the period the teacher already paid for ends, not the
	// day they change plan: cutting immediately takes away something they
	// bought, and the students who lose access did not make the decision.
	//
	// Until then the teacher is told how many will go and may choose which stay.
	DowngradeUsecase interface {
		// Preview says who would be deactivated right now, so the teacher can
		// see the consequence before it happens.
		Preview(ctx context.Context, teacherID string) (*DowngradeOutput, apperrors.ApplicationError)
		// Apply deactivates the excess. `keep` is the teacher's choice; empty
		// means the automatic order.
		Apply(ctx context.Context, teacherID string, keep []string) (*DowngradeOutput, apperrors.ApplicationError)
		// Reactivate brings one student back, refused when the plan is full.
		Reactivate(ctx context.Context, teacherID, studentID string) apperrors.ApplicationError
	}

	downgradeUsecase struct {
		contextFactory appcontext.Factory
	}

	DowngradeData struct {
		MaxStudents int `json:"max_students"`
		// Deactivated is who lost access, or would with Preview.
		Deactivated []string `json:"deactivated"`
	}

	DowngradeOutput struct {
		Data DowngradeData `json:"data"`
	}
)

func NewDowngradeUsecase(contextFactory appcontext.Factory) DowngradeUsecase {
	return &downgradeUsecase{contextFactory: contextFactory}
}

func (u *downgradeUsecase) Preview(ctx context.Context, teacherID string) (*DowngradeOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	scope, appErr := scopeFor(ctx, app, "", teacherID)
	if appErr != nil {
		return nil, appErr
	}
	// Nothing is deactivated for a teacher with no cap, nor on a payments
	// outage: an unreadable plan must not cost anyone their students.
	if !scope.Enforced() {
		return &DowngradeOutput{}, nil
	}

	byActivity, err := app.Repositories.School.ListStudentsByActivity(ctx, scope.SchoolID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}

	return &DowngradeOutput{Data: DowngradeData{
		MaxStudents: scope.Plan.MaxStudents,
		Deactivated: domain.StudentsToDeactivate(byActivity, scope.Plan.MaxStudents, nil),
	}}, nil
}

func (u *downgradeUsecase) Apply(ctx context.Context, teacherID string, keep []string) (*DowngradeOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	scope, appErr := scopeFor(ctx, app, "", teacherID)
	if appErr != nil {
		return nil, appErr
	}
	if !scope.Enforced() {
		return &DowngradeOutput{}, nil
	}

	byActivity, err := app.Repositories.School.ListStudentsByActivity(ctx, scope.SchoolID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}

	deactivated := domain.StudentsToDeactivate(byActivity, scope.Plan.MaxStudents, keep)
	for _, studentID := range deactivated {
		if err := app.Repositories.School.SetMemberActive(ctx, scope.SchoolID, studentID, false); err != nil {
			return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
		}
	}

	return &DowngradeOutput{Data: DowngradeData{
		MaxStudents: scope.Plan.MaxStudents,
		Deactivated: deactivated,
	}}, nil
}

func (u *downgradeUsecase) Reactivate(ctx context.Context, teacherID, studentID string) apperrors.ApplicationError {
	app := u.contextFactory()

	scope, appErr := scopeFor(ctx, app, "", teacherID)
	if appErr != nil {
		return appErr
	}
	if scope.SchoolID == "" {
		return apperrors.NewNotFoundError("no school to reactivate in")
	}

	// Checked against the same rule an addition goes through: reactivating is
	// adding a student back, and the plan does not care how they got there. An
	// outage skips the check for the same reason adding one does.
	if scope.Enforced() {
		used, appErr := studentsUsed(ctx, app, scope.SchoolID)
		if appErr != nil {
			return appErr
		}
		if !(domain.TeacherSubscription{Plan: scope.Plan, StudentsUsed: used}).CanAddStudent() {
			return apperrors.NewApplicationError(mappings.StudentLimitReachedError, nil)
		}
	}

	if err := app.Repositories.School.SetMemberActive(ctx, scope.SchoolID, studentID, true); err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	return nil
}
