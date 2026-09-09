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

	schoolID, plan, appErr := planFor(ctx, app, teacherID)
	if appErr != nil {
		return nil, appErr
	}
	if schoolID == "" {
		return &DowngradeOutput{}, nil
	}

	byActivity, err := app.Repositories.School.ListStudentsByActivity(ctx, schoolID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}

	return &DowngradeOutput{Data: DowngradeData{
		MaxStudents: plan.MaxStudents,
		Deactivated: domain.StudentsToDeactivate(byActivity, plan.MaxStudents, nil),
	}}, nil
}

func (u *downgradeUsecase) Apply(ctx context.Context, teacherID string, keep []string) (*DowngradeOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	schoolID, plan, appErr := planFor(ctx, app, teacherID)
	if appErr != nil {
		return nil, appErr
	}
	if schoolID == "" {
		return &DowngradeOutput{}, nil
	}

	byActivity, err := app.Repositories.School.ListStudentsByActivity(ctx, schoolID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}

	deactivated := domain.StudentsToDeactivate(byActivity, plan.MaxStudents, keep)
	for _, studentID := range deactivated {
		if err := app.Repositories.School.SetMemberActive(ctx, schoolID, studentID, false); err != nil {
			return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
		}
	}

	return &DowngradeOutput{Data: DowngradeData{
		MaxStudents: plan.MaxStudents,
		Deactivated: deactivated,
	}}, nil
}

func (u *downgradeUsecase) Reactivate(ctx context.Context, teacherID, studentID string) apperrors.ApplicationError {
	app := u.contextFactory()

	schoolID, plan, appErr := planFor(ctx, app, teacherID)
	if appErr != nil {
		return appErr
	}
	if schoolID == "" {
		return apperrors.NewNotFoundError("no school to reactivate in")
	}

	used, err := app.Repositories.School.CountStudents(ctx, schoolID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	// Checked against the same rule an addition goes through: reactivating is
	// adding a student back, and the plan does not care how they got there.
	if !(domain.TeacherSubscription{Plan: plan, StudentsUsed: used}).CanAddStudent() {
		return apperrors.NewApplicationError(mappings.StudentLimitReachedError, nil)
	}

	if err := app.Repositories.School.SetMemberActive(ctx, schoolID, studentID, true); err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	return nil
}

// planFor resolves the school a teacher's plan applies to, and what it allows.
// An empty school id means there is nothing to enforce: no school of their own,
// or one invoiced outside the product.
func planFor(ctx context.Context, app *appcontext.Context, teacherID string) (string, domain.TeacherPlan, apperrors.ApplicationError) {
	school, err := app.Repositories.School.GetPersonal(ctx, teacherID)
	if err != nil {
		return "", domain.TeacherPlan{}, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	if school == nil || school.Billing == domain.SchoolBillingDirect {
		return "", domain.TeacherPlan{}, nil
	}

	plan := domain.FreePlan
	entitlement, err := app.Integrations.Payments.GetEntitlement(ctx, teacherID)
	if err != nil {
		// Same direction as everywhere else: an unreadable plan must not cost
		// anyone their students. Nothing is deactivated on a payments outage.
		return "", domain.TeacherPlan{}, nil
	}
	if entitlement != nil && entitlement.Active {
		plan = domain.PlanFromMetadata(0, "", entitlement.Metadata)
	}
	return school.ID, plan, nil
}
