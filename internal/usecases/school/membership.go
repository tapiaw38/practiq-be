package school

import (
	"context"
	"log"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

func JoinTeacherSchool(ctx context.Context, app *appcontext.Context, teacherID, studentID string) {
	JoinSchool(ctx, app, "", teacherID, studentID)
}

func JoinSchool(ctx context.Context, app *appcontext.Context, schoolID, teacherID, studentID string) {
	if schoolID == "" {
		school, err := app.Repositories.School.GetPersonal(ctx, teacherID)
		if err != nil {
			log.Printf("[schools] school lookup failed teacher_id=%s err=%v", teacherID, err)
			return
		}
		if school == nil {
			return
		}
		schoolID = school.ID
	}
	if err := RequireActive(ctx, app, schoolID); err != nil {
		log.Printf("[schools] membership skipped inactive school_id=%s", schoolID)
		return
	}

	if err := app.Repositories.School.AddMember(ctx, domain.SchoolMember{
		SchoolID: schoolID,
		UserID:   studentID,
		Role:     domain.SchoolRoleStudent,
	}); err != nil {
		log.Printf("[schools] membership failed school_id=%s user_id=%s err=%v", schoolID, studentID, err)
	}
}
