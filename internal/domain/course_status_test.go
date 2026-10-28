package domain

import "testing"

func TestCourseStatusRules(t *testing.T) {
	cases := []struct {
		status          string
		wantVisible     bool
		wantAcceptsWork bool
		wantEnrollment  bool
	}{
		{CourseStatusDraft, false, false, false},
		{CourseStatusPublished, true, true, true},
		{CourseStatusArchived, true, false, false},

		{"", false, false, false},
	}

	for _, tc := range cases {
		t.Run(tc.status, func(t *testing.T) {
			c := Course{Status: tc.status}
			if got := c.VisibleToStudents(); got != tc.wantVisible {
				t.Errorf("VisibleToStudents = %v, want %v", got, tc.wantVisible)
			}
			if got := c.AcceptsWork(); got != tc.wantAcceptsWork {
				t.Errorf("AcceptsWork = %v, want %v", got, tc.wantAcceptsWork)
			}
			if got := c.AcceptsEnrollment(); got != tc.wantEnrollment {
				t.Errorf("AcceptsEnrollment = %v, want %v", got, tc.wantEnrollment)
			}
		})
	}
}
