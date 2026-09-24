package school

import (
	"context"
	"testing"

	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories"
	schoolRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/school"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

// fakeSchoolsByUser keys membership by user, unlike fakeSchools in
// scope_test.go, because these checks compare different users' memberships
// against each other.
type fakeSchoolsByUser struct {
	schoolRepo.Repository
	byUser map[string][]domain.SchoolMember
}

func (f *fakeSchoolsByUser) ListForUser(_ context.Context, userID string) ([]domain.SchoolMember, error) {
	return f.byUser[userID], nil
}

func appWithUsers(byUser map[string][]domain.SchoolMember) *appcontext.Context {
	return &appcontext.Context{
		Repositories: &repositories.Repositories{School: &fakeSchoolsByUser{byUser: byUser}},
	}
}

func TestEnsureCanLinkTeacherStudent(t *testing.T) {
	admin := []domain.SchoolMember{{SchoolID: "school-a", Role: domain.SchoolRoleAdmin, Active: true}}

	cases := []struct {
		name         string
		byUser       map[string][]domain.SchoolMember
		isSuperAdmin bool
		wantAllowed  bool
	}{
		{
			name: "an admin may link a teacher and student who share their school",
			byUser: map[string][]domain.SchoolMember{
				"admin":   admin,
				"teacher": {{SchoolID: "school-a", Role: domain.SchoolRoleTeacher, Active: true}},
				"student": {{SchoolID: "school-a", Role: domain.SchoolRoleStudent, Active: true}},
			},
			wantAllowed: true,
		},
		{
			// The reason the check exists. Without it, an admin of school A
			// could link a teacher from school B to a student from school C.
			name: "an admin may not link people from a school they don't administer",
			byUser: map[string][]domain.SchoolMember{
				"admin":   admin,
				"teacher": {{SchoolID: "school-b", Role: domain.SchoolRoleTeacher, Active: true}},
				"student": {{SchoolID: "school-c", Role: domain.SchoolRoleStudent, Active: true}},
			},
		},
		{
			name: "the teacher must actually belong to the admin's school",
			byUser: map[string][]domain.SchoolMember{
				"admin":   admin,
				"teacher": {{SchoolID: "school-b", Role: domain.SchoolRoleTeacher, Active: true}},
				"student": {{SchoolID: "school-a", Role: domain.SchoolRoleStudent, Active: true}},
			},
		},
		{
			name: "a deactivated admin membership does not count",
			byUser: map[string][]domain.SchoolMember{
				"admin":   {{SchoolID: "school-a", Role: domain.SchoolRoleAdmin, Active: false}},
				"teacher": {{SchoolID: "school-a", Role: domain.SchoolRoleTeacher, Active: true}},
				"student": {{SchoolID: "school-a", Role: domain.SchoolRoleStudent, Active: true}},
			},
		},
		{
			name:         "a superadmin passes regardless of membership",
			isSuperAdmin: true,
			wantAllowed:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			appErr := EnsureCanLinkTeacherStudent(
				context.Background(), appWithUsers(tc.byUser), "admin", tc.isSuperAdmin, "teacher", "student",
			)
			if allowed := appErr == nil; allowed != tc.wantAllowed {
				t.Fatalf("allowed = %v, want %v", allowed, tc.wantAllowed)
			}
		})
	}
}

func TestEnsureCanViewAssignmentsFor(t *testing.T) {
	admin := []domain.SchoolMember{{SchoolID: "school-a", Role: domain.SchoolRoleAdmin, Active: true}}

	cases := []struct {
		name         string
		requesterID  string
		targetID     string
		byUser       map[string][]domain.SchoolMember
		isSuperAdmin bool
		wantAllowed  bool
	}{
		{
			name:        "a user may always view their own assignments",
			requesterID: "teacher-1",
			targetID:    "teacher-1",
			wantAllowed: true,
		},
		{
			name:        "an admin may view a member of their school",
			requesterID: "admin",
			targetID:    "teacher",
			byUser: map[string][]domain.SchoolMember{
				"admin":   admin,
				"teacher": {{SchoolID: "school-a", Role: domain.SchoolRoleTeacher, Active: true}},
			},
			wantAllowed: true,
		},
		{
			name:        "an admin may not view someone outside their school",
			requesterID: "admin",
			targetID:    "teacher",
			byUser: map[string][]domain.SchoolMember{
				"admin":   admin,
				"teacher": {{SchoolID: "school-b", Role: domain.SchoolRoleTeacher, Active: true}},
			},
		},
		{
			name:         "a superadmin passes regardless of membership",
			requesterID:  "admin",
			targetID:     "anyone",
			isSuperAdmin: true,
			wantAllowed:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			appErr := EnsureCanViewAssignmentsFor(
				context.Background(), appWithUsers(tc.byUser), tc.requesterID, tc.isSuperAdmin, tc.targetID,
			)
			if allowed := appErr == nil; allowed != tc.wantAllowed {
				t.Fatalf("allowed = %v, want %v", allowed, tc.wantAllowed)
			}
		})
	}
}
