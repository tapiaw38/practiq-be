package school

import (
	"context"
	"testing"

	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories"
	courseRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/payments"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

type fakeCourses struct {
	courseRepo.Repository
	course *domain.Course
}

func (f *fakeCourses) Get(context.Context, string) (*domain.Course, error) {
	return f.course, nil
}

// payingTeacher keeps these tests about membership. Whether the school is paid
// up is its own rule, with its own tests; leaving it unset would panic here and
// hide what each case is actually asserting.
type payingTeacher struct{ payments.Client }

func (payingTeacher) GetEntitlement(context.Context, string) (*payments.Entitlement, error) {
	return &payments.Entitlement{Active: true, Metadata: map[string]any{"max_students": float64(30)}}, nil
}

func appWithCourse(course *domain.Course, members ...domain.SchoolMember) *appcontext.Context {
	return &appcontext.Context{
		Repositories: &repositories.Repositories{
			School: &fakeSchools{members: members},
			Course: &fakeCourses{course: course},
		},
		Integrations: &integrations.Integrations{Payments: payingTeacher{}},
	}
}

func TestEnsureCanManageCourse(t *testing.T) {
	course := &domain.Course{ID: "course-1", TeacherID: "teacher-1", SchoolID: "school-1"}

	cases := []struct {
		name         string
		course       *domain.Course
		members      []domain.SchoolMember
		isSuperAdmin bool
		wantAllowed  bool
	}{
		{
			name:        "the teacher who owns it",
			course:      course,
			wantAllowed: true,
		},
		{
			name:        "an admin of the school it belongs to",
			course:      &domain.Course{ID: "course-1", TeacherID: "someone-else", SchoolID: "school-1"},
			members:     []domain.SchoolMember{{SchoolID: "school-1", Role: domain.SchoolRoleAdmin, Active: true}},
			wantAllowed: true,
		},
		{
			name:         "a platform superadmin",
			course:       &domain.Course{ID: "course-1", TeacherID: "someone-else", SchoolID: "school-1"},
			isSuperAdmin: true,
			wantAllowed:  true,
		},
		{
			name:    "another teacher of the same school",
			course:  &domain.Course{ID: "course-1", TeacherID: "someone-else", SchoolID: "school-1"},
			members: []domain.SchoolMember{{SchoolID: "school-1", Role: domain.SchoolRoleTeacher, Active: true}},
		},
		{
			name:    "an admin of a different school",
			course:  &domain.Course{ID: "course-1", TeacherID: "someone-else", SchoolID: "school-1"},
			members: []domain.SchoolMember{{SchoolID: "school-2", Role: domain.SchoolRoleAdmin, Active: true}},
		},
		{
			// Course.Get joins an active school, so a closed one reads as no
			// course at all — for the operator too.
			name:         "a course that does not read back",
			course:       nil,
			isSuperAdmin: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, appErr := EnsureCanManageCourse(
				context.Background(), appWithCourse(tc.course, tc.members...),
				"teacher-1", tc.isSuperAdmin, "course-1",
			)
			if allowed := appErr == nil; allowed != tc.wantAllowed {
				t.Fatalf("allowed = %v, want %v", allowed, tc.wantAllowed)
			}
		})
	}
}

func TestOwnedSchoolIDSelected(t *testing.T) {
	cases := []struct {
		name         string
		members      []domain.SchoolMember
		isSuperAdmin bool
		selected     string
		want         string
	}{
		{
			name:     "an admin gets the school they administer",
			members:  []domain.SchoolMember{{SchoolID: "school-1", Role: domain.SchoolRoleAdmin, Active: true}},
			selected: "school-1",
			want:     "school-1",
		},
		{
			name:     "an admin may not file under another school",
			members:  []domain.SchoolMember{{SchoolID: "school-1", Role: domain.SchoolRoleAdmin, Active: true}},
			selected: "school-2",
		},
		{
			// Without this a superadmin could create a course in an
			// institution but not the grade or subject that course needs.
			name:         "a superadmin gets the school they selected",
			isSuperAdmin: true,
			selected:     "school-2",
			want:         "school-2",
		},
		{
			name:         "a superadmin who selected nothing is refused",
			isSuperAdmin: true,
			selected:     "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, appErr := OwnedSchoolIDSelected(
				context.Background(), appWith(tc.members...), "user-1", tc.isSuperAdmin, tc.selected,
			)
			if tc.want == "" {
				if appErr == nil {
					t.Fatalf("got %q, want refusal", got)
				}
				return
			}
			if appErr != nil {
				t.Fatalf("unexpected refusal: %v", appErr)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
