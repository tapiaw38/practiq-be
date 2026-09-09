package school

import (
	"context"
	"testing"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func TestEnsureAdministers(t *testing.T) {
	cases := []struct {
		name         string
		members      []domain.SchoolMember
		isSuperAdmin bool
		schoolID     string
		wantAllowed  bool
	}{
		{
			name:        "an admin may touch their own school",
			members:     []domain.SchoolMember{{SchoolID: "mine", Role: domain.SchoolRoleAdmin, Active: true}},
			schoolID:    "mine",
			wantAllowed: true,
		},
		{
			// The reason the guard exists. Opening these routes beyond the
			// platform superadmin without it lets any teacher rename another
			// school's grades.
			name:     "an admin may not touch another school",
			members:  []domain.SchoolMember{{SchoolID: "mine", Role: domain.SchoolRoleAdmin, Active: true}},
			schoolID: "someone-elses",
		},
		{
			name:     "a teacher of the school is not its admin",
			members:  []domain.SchoolMember{{SchoolID: "school", Role: domain.SchoolRoleTeacher, Active: true}},
			schoolID: "school",
		},
		{
			name:     "a deactivated admin is not an admin",
			members:  []domain.SchoolMember{{SchoolID: "school", Role: domain.SchoolRoleAdmin, Active: false}},
			schoolID: "school",
		},
		{
			// Rows from before the split, or created while nothing assigned a
			// school. The first teacher to find one must not take it over.
			name:     "a row with no school belongs to nobody",
			members:  []domain.SchoolMember{{SchoolID: "school", Role: domain.SchoolRoleAdmin, Active: true}},
			schoolID: "",
		},
		{
			name:         "a superadmin passes anywhere",
			isSuperAdmin: true,
			schoolID:     "any",
			wantAllowed:  true,
		},
		{
			name:         "a superadmin passes even on an ownerless row",
			isSuperAdmin: true,
			schoolID:     "",
			wantAllowed:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			appErr := EnsureAdministers(
				context.Background(), appWith(tc.members...), "user-1", tc.isSuperAdmin, tc.schoolID,
			)
			if allowed := appErr == nil; allowed != tc.wantAllowed {
				t.Fatalf("allowed = %v, want %v", allowed, tc.wantAllowed)
			}
		})
	}
}
