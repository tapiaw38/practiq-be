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
	DowngradeUsecase interface {
		Preview(ctx context.Context, teacherID, bearerToken string) (*DowngradeOutput, apperrors.ApplicationError)

		Apply(ctx context.Context, teacherID string, keep []string) (*DowngradeOutput, apperrors.ApplicationError)

		Reactivate(ctx context.Context, teacherID, studentID string) apperrors.ApplicationError
	}

	downgradeUsecase struct {
		contextFactory appcontext.Factory
	}

	DowngradeStudent struct {
		ID   string `json:"id"`
		Name string `json:"name"`

		LastPracticedAt string `json:"last_practiced_at,omitempty"`

		Keeps bool `json:"keeps"`
	}

	DowngradeData struct {
		MaxStudents int `json:"max_students"`

		Deactivated []string `json:"deactivated"`

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

	if !scope.Enforced() {

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
