package subscription

import (
	"context"
	"errors"
	"testing"

	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories"
	teacherstudentassignment "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/teacher_student_assignment"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/payments"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

// fakeAssignments answers the two questions the allowance check asks.
type fakeAssignments struct {
	teacherstudentassignment.Repository
	hasAccess bool
	count     int
}

func (f *fakeAssignments) HasAccess(context.Context, string, string) (bool, error) {
	return f.hasAccess, nil
}

func (f *fakeAssignments) CountStudents(context.Context, string) (int, error) {
	return f.count, nil
}

type fakePayments struct {
	payments.Client
	entitlement *payments.Entitlement
	err         error
}

func (f *fakePayments) GetEntitlement(context.Context, string) (*payments.Entitlement, error) {
	return f.entitlement, f.err
}

func appWith(assignments *fakeAssignments, pay *fakePayments) *appcontext.Context {
	return &appcontext.Context{
		Repositories: &repositories.Repositories{TeacherStudentAssignment: assignments},
		Integrations: &integrations.Integrations{Payments: pay},
	}
}

func paidPlan(maxStudents int) *payments.Entitlement {
	return &payments.Entitlement{
		Active:   true,
		Metadata: map[string]any{"max_students": float64(maxStudents)},
	}
}

func TestEnsureCanAddStudent(t *testing.T) {
	cases := []struct {
		name        string
		assignments *fakeAssignments
		payments    *fakePayments
		wantRefused bool
	}{
		{
			name:        "there is room on the plan",
			assignments: &fakeAssignments{count: 3},
			payments:    &fakePayments{entitlement: paidPlan(5)},
		},
		{
			name:        "the plan is full",
			assignments: &fakeAssignments{count: 5},
			payments:    &fakePayments{entitlement: paidPlan(5)},
			wantRefused: true,
		},
		{
			// These paths are idempotent. Re-running one must not start failing
			// because the plan filled up in between; nothing is being added.
			name:        "a student the teacher already has is always allowed",
			assignments: &fakeAssignments{hasAccess: true, count: 99},
			payments:    &fakePayments{entitlement: paidPlan(5)},
		},
		{
			name:        "no subscription means the free allowance",
			assignments: &fakeAssignments{count: domain.FreePlan.MaxStudents},
			payments:    &fakePayments{entitlement: &payments.Entitlement{Active: false}},
			wantRefused: true,
		},
		{
			// The dangerous default. Treating an unreadable plan as the free
			// one would stop every paying teacher from working whenever the
			// payments service hiccups; one extra student costs far less.
			name:        "a payments outage does not block the teacher",
			assignments: &fakeAssignments{count: 500},
			payments:    &fakePayments{err: errors.New("payments is down")},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := appWith(tc.assignments, tc.payments)

			appErr := EnsureCanAddStudent(context.Background(), app, "teacher-1", "student-1")

			if tc.wantRefused && appErr == nil {
				t.Fatal("expected the link to be refused")
			}
			if !tc.wantRefused && appErr != nil {
				t.Fatalf("expected the link to be allowed, got %v", appErr)
			}
		})
	}
}
