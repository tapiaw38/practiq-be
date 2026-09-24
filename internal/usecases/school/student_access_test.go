package school

import (
	"testing"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

// The teacher is told a downgrade makes these students lose access. Until this
// guard existed it did not: deactivating only stopped them counting against
// the plan, so ten students dropped from a plan of fifteen kept working.
func TestEnsureStudentCanWork(t *testing.T) {
	course := &domain.Course{ID: "course-1", TeacherID: "teacher-1", SchoolID: "school-1"}

	cases := []struct {
		name        string
		course      *domain.Course
		members     []domain.SchoolMember
		wantAllowed bool
	}{
		{
			name:        "an active student works",
			course:      course,
			members:     []domain.SchoolMember{{SchoolID: "school-1", UserID: "student-1", Role: "student", Active: true}},
			wantAllowed: true,
		},
		{
			name:        "a deactivated student does not",
			course:      course,
			members:     []domain.SchoolMember{{SchoolID: "school-1", UserID: "student-1", Role: "student", Active: false}},
			wantAllowed: false,
		},
		{
			name:   "deactivated somewhere else is not this school's business",
			course: course,
			members: []domain.SchoolMember{
				{SchoolID: "school-2", UserID: "student-1", Role: "student", Active: false},
				{SchoolID: "school-1", UserID: "student-1", Role: "student", Active: true},
			},
			wantAllowed: true,
		},
		{
			name:        "no membership at all is not somebody's downgrade",
			course:      course,
			members:     nil,
			wantAllowed: true,
		},
		{
			name:        "a course outside any school cannot be judged",
			course:      &domain.Course{ID: "course-1", TeacherID: "teacher-1"},
			members:     []domain.SchoolMember{{SchoolID: "school-1", UserID: "student-1", Role: "student", Active: false}},
			wantAllowed: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := appWithCourse(tc.course, tc.members...)
			err := EnsureStudentCanWork(t.Context(), app, "student-1", "course-1")
			if (err == nil) != tc.wantAllowed {
				t.Fatalf("allowed = %v, want %v (err=%v)", err == nil, tc.wantAllowed, err)
			}
		})
	}
}
