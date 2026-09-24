package subscription

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
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
		Preview(ctx context.Context, teacherID, bearerToken string) (*DowngradeOutput, apperrors.ApplicationError)
		// Apply deactivates the excess. `keep` is the teacher's choice; empty
		// means the automatic order.
		Apply(ctx context.Context, teacherID string, keep []string) (*DowngradeOutput, apperrors.ApplicationError)
		// Reactivate brings one student back, refused when the plan is full.
		Reactivate(ctx context.Context, teacherID, studentID string) apperrors.ApplicationError
	}

	downgradeUsecase struct {
		contextFactory appcontext.Factory
	}

	// DowngradeStudent is one candidate, named and dated so the teacher can
	// actually choose. A list of opaque ids is not a choice.
	DowngradeStudent struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		// LastPracticedAt is empty for somebody who never practised — the
		// reason they are first in line.
		LastPracticedAt string `json:"last_practiced_at,omitempty"`
		// Keeps is what the automatic order would do, so the screen can start
		// from it rather than from an empty form.
		Keeps bool `json:"keeps"`
	}

	DowngradeData struct {
		MaxStudents int `json:"max_students"`
		// Deactivated is who lost access, or would with Preview.
		Deactivated []string `json:"deactivated"`
		// Students is everyone the cap applies to, in the order it applies
		// them. Only present on Preview, and only when there is a choice.
		Students []DowngradeStudent `json:"students,omitempty"`
	}

	DowngradeOutput struct {
		Data DowngradeData `json:"data"`
	}
)

func NewDowngradeUsecase(contextFactory appcontext.Factory) DowngradeUsecase {
	return &downgradeUsecase{contextFactory: contextFactory}
}

func (u *downgradeUsecase) Preview(ctx context.Context, teacherID, bearerToken string) (*DowngradeOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	scope, appErr := scopeFor(ctx, app, "", teacherID)
	if appErr != nil {
		return nil, appErr
	}
	// Nothing is deactivated for a teacher with no cap, nor on a payments
	// outage: an unreadable plan must not cost anyone their students.
	if !scope.Enforced() {
		// JSON encodes a nil slice as null. This is a collection in the public
		// contract, so keep it an empty array even when no plan is enforced.
		return &DowngradeOutput{Data: DowngradeData{Deactivated: []string{}}}, nil
	}

	activity, err := app.Repositories.School.ListStudentsWithActivity(ctx, scope.SchoolID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}

	byActivity := make([]string, 0, len(activity))
	for _, student := range activity {
		byActivity = append(byActivity, student.UserID)
	}
	deactivated := domain.StudentsToDeactivate(byActivity, scope.Plan.MaxStudents, nil)
	if deactivated == nil {
		deactivated = []string{}
	}

	data := DowngradeData{MaxStudents: scope.Plan.MaxStudents, Deactivated: deactivated}
	if len(deactivated) > 0 {
		data.Students = describeCandidates(ctx, app, bearerToken, activity, deactivated)
	}
	return &DowngradeOutput{Data: data}, nil
}

// describeCandidates names and dates everyone the cap applies to.
//
// Names come from auth-api and are best effort: a teacher choosing between ten
// students is better served by nine names and one id than by an error page.
func describeCandidates(
	ctx context.Context,
	app *appcontext.Context,
	bearerToken string,
	activity []domain.StudentActivity,
	deactivated []string,
) []DowngradeStudent {
	ids := make([]string, 0, len(activity))
	for _, student := range activity {
		ids = append(ids, student.UserID)
	}
	names, appErr := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, ids)
	if appErr != nil {
		log.Printf("[subscription] student names unavailable err=%v", appErr)
		names = nil
	}

	losing := make(map[string]bool, len(deactivated))
	for _, id := range deactivated {
		losing[id] = true
	}

	out := make([]DowngradeStudent, 0, len(activity))
	for _, student := range activity {
		candidate := DowngradeStudent{ID: student.UserID, Keeps: !losing[student.UserID]}
		if info, ok := names[student.UserID]; ok {
			candidate.Name = strings.TrimSpace(info.FirstName + " " + info.LastName)
		}
		if candidate.Name == "" {
			candidate.Name = "Alumno sin nombre"
		}
		if student.LastPracticed != nil {
			candidate.LastPracticedAt = student.LastPracticed.UTC().Format(time.RFC3339)
		}
		out = append(out, candidate)
	}
	return out
}

func (u *downgradeUsecase) Apply(ctx context.Context, teacherID string, keep []string) (*DowngradeOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	scope, appErr := scopeFor(ctx, app, "", teacherID)
	if appErr != nil {
		return nil, appErr
	}
	if !scope.Enforced() {
		return &DowngradeOutput{Data: DowngradeData{Deactivated: []string{}}}, nil
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
