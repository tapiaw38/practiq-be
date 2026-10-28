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
	DowngradePreviewUsecase interface {
		Execute(ctx context.Context, teacherID, bearerToken string) (*DowngradePreviewOutput, apperrors.ApplicationError)
	}

	downgradePreviewUsecase struct {
		contextFactory appcontext.Factory
	}

	DowngradePreviewOutput struct {
		Data DowngradePreviewData `json:"data"`
	}

	DowngradePreviewData struct {
		MaxStudents int                `json:"max_students"`
		Deactivated []string           `json:"deactivated"`
		Students    []DowngradeStudent `json:"students,omitempty"`
	}

	DowngradeStudent struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		LastPracticedAt string `json:"last_practiced_at,omitempty"`
		Keeps           bool   `json:"keeps"`
	}
)

func NewDowngradePreviewUsecase(contextFactory appcontext.Factory) DowngradePreviewUsecase {
	return &downgradePreviewUsecase{contextFactory: contextFactory}
}

func (u *downgradePreviewUsecase) Execute(ctx context.Context, teacherID, bearerToken string) (*DowngradePreviewOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	scope, appErr := scopeFor(ctx, app, "", teacherID)
	if appErr != nil {
		return nil, appErr
	}
	if !scope.Enforced() {
		return &DowngradePreviewOutput{Data: DowngradePreviewData{Deactivated: []string{}}}, nil
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
	data := DowngradePreviewData{MaxStudents: scope.Plan.MaxStudents, Deactivated: deactivated}
	if len(deactivated) > 0 {
		data.Students = describeDowngradeCandidates(ctx, app, bearerToken, activity, deactivated)
	}
	return &DowngradePreviewOutput{Data: data}, nil
}

func describeDowngradeCandidates(ctx context.Context, app *appcontext.Context, bearerToken string, activity []domain.StudentActivity, deactivated []string) []DowngradeStudent {
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
