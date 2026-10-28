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
	teacherID string
	studentID string
	schoolID  string
}

func (s *listUsecaseSpy) Execute(_ context.Context, teacherID, studentID, schoolID string) (*ucCourse.ListOutput, apperrors.ApplicationError) {
	s.teacherID, s.studentID, s.schoolID = teacherID, studentID, schoolID
	return &ucCourse.ListOutput{Data: []ucCourse.CourseData{}}, nil
}

func listAs(t *testing.T, userID string, roles []auth.RoleClaim, schoolID string) *listUsecaseSpy {
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
	return spy
}

func TestListNarrowsATeacherToTheirOwnCourses(t *testing.T) {
	got := listAs(t, "teacher-1", []auth.RoleClaim{{Name: "user"}}, "school-1")

	if got.teacherID != "teacher-1" {
		t.Fatalf("TeacherID = %q, want the caller", got.teacherID)
	}
	if got.schoolID != "school-1" {
		t.Fatalf("SchoolID = %q, want the header value", got.schoolID)
	}
}

func TestListShowsASuperAdminEveryCourseOfTheSelectedSchool(t *testing.T) {
	got := listAs(t, "operator-1", []auth.RoleClaim{{Name: "superadmin"}}, "school-1")

	if got.teacherID != "" {
		t.Fatalf("TeacherID = %q, want no ownership filter for a superadmin", got.teacherID)
	}
	if got.studentID != "" {
		t.Fatalf("StudentID = %q, want no student filter on role=teacher", got.studentID)
	}
	if got.schoolID != "school-1" {
		t.Fatalf("SchoolID = %q, want the answer still bounded by the school", got.schoolID)
	}
}
