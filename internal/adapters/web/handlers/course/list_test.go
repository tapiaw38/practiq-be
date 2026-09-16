package course

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/platform/auth"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	ucCourse "github.com/tapiaw38/practiq-be/internal/usecases/course"
)

type listUsecaseSpy struct {
	got ucCourse.ListInput
}

func (s *listUsecaseSpy) Execute(_ context.Context, input ucCourse.ListInput) (*ucCourse.ListOutput, apperrors.ApplicationError) {
	s.got = input
	return &ucCourse.ListOutput{Data: []ucCourse.CourseData{}}, nil
}

// listAs runs the handler with the identity a request would carry after the
// auth middleware, and reports the filter the usecase was asked for.
func listAs(t *testing.T, userID string, roles []auth.RoleClaim, schoolID string) ucCourse.ListInput {
	t.Helper()
	gin.SetMode(gin.TestMode)

	spy := &listUsecaseSpy{}
	app := gin.New()
	app.GET("/api/courses", func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("roles", roles)
		NewListHandler(spy)(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/courses?role=teacher", nil)
	if schoolID != "" {
		req.Header.Set("X-School-ID", schoolID)
	}
	res := httptest.NewRecorder()
	app.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.Code)
	}
	return spy.got
}

func TestListNarrowsATeacherToTheirOwnCourses(t *testing.T) {
	got := listAs(t, "teacher-1", []auth.RoleClaim{{Name: "user"}}, "school-1")

	if got.TeacherID != "teacher-1" {
		t.Fatalf("TeacherID = %q, want the caller", got.TeacherID)
	}
	if got.SchoolID != "school-1" {
		t.Fatalf("SchoolID = %q, want the header value", got.SchoolID)
	}
}

// A superadmin opening a school operates it rather than teaching in it.
// Filtering by ownership showed them an empty school and had them re-create
// courses that were already there.
func TestListShowsASuperAdminEveryCourseOfTheSelectedSchool(t *testing.T) {
	got := listAs(t, "operator-1", []auth.RoleClaim{{Name: "superadmin"}}, "school-1")

	if got.TeacherID != "" {
		t.Fatalf("TeacherID = %q, want no ownership filter for a superadmin", got.TeacherID)
	}
	if got.StudentID != "" {
		t.Fatalf("StudentID = %q, want no student filter on role=teacher", got.StudentID)
	}
	if got.SchoolID != "school-1" {
		t.Fatalf("SchoolID = %q, want the answer still bounded by the school", got.SchoolID)
	}
}
