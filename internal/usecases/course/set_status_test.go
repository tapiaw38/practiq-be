package course

import (
	"testing"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func TestValidCourseStatus(t *testing.T) {
	cases := []struct {
		status string
		want   bool
	}{
		{domain.CourseStatusDraft, true},
		{domain.CourseStatusPublished, true},
		{domain.CourseStatusArchived, true},
		{"", false},
		{"Published", false},
		{"deleted", false},
		{" published", false},
	}

	for _, tc := range cases {
		t.Run(tc.status, func(t *testing.T) {
			if got := ValidCourseStatus(tc.status); got != tc.want {
				t.Fatalf("ValidCourseStatus(%q) = %v, want %v", tc.status, got, tc.want)
			}
		})
	}
}
