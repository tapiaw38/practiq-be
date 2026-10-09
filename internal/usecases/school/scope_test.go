package school

import (
	"context"
	"testing"

	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories"
	schoolRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/school"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

type fakeSchools struct {
	schoolRepo.Repository
	members            []domain.SchoolMember
	studentsByActivity []string
}

func (f *fakeSchools) ListForUser(context.Context, string) ([]domain.SchoolMember, error) {
	return f.members, nil
}

func (f *fakeSchools) ListStudentsByActivity(context.Context, string) ([]string, error) {
	return f.studentsByActivity, nil
}

func (f *fakeSchools) Get(_ context.Context, id string) (*domain.School, error) {
	return &domain.School{ID: id, Status: domain.SchoolStatusActive}, nil
}

func appWith(members ...domain.SchoolMember) *appcontext.Context {
	return &appcontext.Context{
		Repositories: &repositories.Repositories{School: &fakeSchools{members: members}},
	}
}

func TestScopeFor(t *testing.T) {
	cases := []struct {
		name         string
		members      []domain.SchoolMember
		isSuperAdmin bool
		wantIDs      []string
		wantAll      bool
		wantEmpty    bool
	}{
		{
			name:         "a superadmin is not narrowed at all",
			isSuperAdmin: true,
			wantAll:      true,
		},
		{

			name:         "a superadmin belonging to nothing still sees everything",
			members:      nil,
			isSuperAdmin: true,
			wantAll:      true,
		},
		{
			name: "a teacher sees the schools they belong to",
			members: []domain.SchoolMember{
				{SchoolID: "a", Role: domain.SchoolRoleAdmin, Active: true},
				{SchoolID: "b", Role: domain.SchoolRoleTeacher, Active: true},
			},
			wantIDs: []string{"a", "b"},
		},
		{

			name: "a deactivated membership does not count",
			members: []domain.SchoolMember{
				{SchoolID: "a", Role: domain.SchoolRoleStudent, Active: false},
			},
			wantEmpty: true,
		},
		{

			name:      "belonging to nothing sees nothing, not everything",
			members:   nil,
			wantEmpty: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			scope, err := ScopeFor(context.Background(), appWith(tc.members...), "user-1", tc.isSuperAdmin)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if scope.Unrestricted != tc.wantAll {
				t.Fatalf("unrestricted = %v, want %v", scope.Unrestricted, tc.wantAll)
			}
			if scope.Empty() != tc.wantEmpty {
				t.Fatalf("empty = %v, want %v", scope.Empty(), tc.wantEmpty)
			}
			if len(scope.SchoolIDs) != len(tc.wantIDs) {
				t.Fatalf("got %v, want %v", scope.SchoolIDs, tc.wantIDs)
			}
			for i, id := range tc.wantIDs {
				if scope.SchoolIDs[i] != id {
					t.Fatalf("got %v, want %v", scope.SchoolIDs, tc.wantIDs)
				}
			}
		})
	}
}

func TestSuperAdminScopePassesNilNotEmpty(t *testing.T) {
	scope, err := ScopeFor(context.Background(), appWith(), "user-1", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scope.SchoolIDs != nil {
		t.Fatalf("superadmin school ids = %v, want nil", scope.SchoolIDs)
	}
}
