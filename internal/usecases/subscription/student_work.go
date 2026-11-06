package subscription

import (
	"context"
	"log"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

func StudentsCanWork(ctx context.Context, app *appcontext.Context, schoolID, ownerID string) bool {
	scope, appErr := scopeFor(ctx, app, schoolID, ownerID)
	if appErr != nil {
		return true
	}
	switch scope.State {
	case capNone, capUnknown:
		return true
	}
	if scope.Active || scope.GraceEndsAt != nil {
		return true
	}

	return scope.Plan.MaxStudents > 0
}

func StudentMayWork(ctx context.Context, app *appcontext.Context, schoolID, ownerID, studentID string) bool {
	if !StudentsCanWork(ctx, app, schoolID, ownerID) {
		return false
	}
	scope, appErr := scopeFor(ctx, app, schoolID, ownerID)
	if appErr != nil || !scope.Enforced() || scope.SchoolID == "" {
		return true
	}

	byActivity, err := app.Repositories.School.ListStudentsByActivity(ctx, scope.SchoolID)
	if err != nil {
		log.Printf("[subscription] activity lookup failed school_id=%s err=%v", scope.SchoolID, err)
		return true
	}
	for _, id := range domain.StudentsToDeactivate(byActivity, scope.Plan.MaxStudents, nil) {
		if id == studentID {
			return false
		}
	}
	return true
}
