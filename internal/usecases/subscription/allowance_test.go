package subscription

import (
	"context"
	"errors"
	"testing"

	"time"

	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories"
	schoolRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/school"
	teacherstudentassignment "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/teacher_student_assignment"
	userprofile "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/user_profile"
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

type fakeSchools struct {
	schoolRepo.Repository
	school   *domain.School
	students int
}

func (f *fakeSchools) GetPersonal(context.Context, string) (*domain.School, error) {
	return f.school, nil
}

func (f *fakeSchools) CountStudents(context.Context, string) (int, error) {
	return f.students, nil
}

type fakePayments struct {
	payments.Client
	entitlement *payments.Entitlement
	err         error
}

func (f *fakePayments) GetEntitlement(context.Context, string) (*payments.Entitlement, error) {
	return f.entitlement, f.err
}

// fakeProfiles carries the one thing the free plan depends on: how long ago
// the teacher signed up, which is what the trial is counted from.
type fakeProfiles struct {
	userprofile.Repository
	createdAt time.Time
}

func (f *fakeProfiles) Get(context.Context, string) (*domain.UserProfile, error) {
	return &domain.UserProfile{CreatedAt: f.createdAt}, nil
}

// withinTrial and pastTrial are the two sides of the free month.
func withinTrial() *fakeProfiles {
	return &fakeProfiles{createdAt: time.Now().AddDate(0, 0, -1)}
}

func pastTrial() *fakeProfiles {
	return &fakeProfiles{createdAt: time.Now().AddDate(0, 0, -(domain.FreePlan.TrialDays + 1))}
}

func appWith(assignments *fakeAssignments, pay *fakePayments, schools *fakeSchools, profiles *fakeProfiles) *appcontext.Context {
	if profiles == nil {
		profiles = withinTrial()
	}
	return &appcontext.Context{
		Repositories: &repositories.Repositories{
			TeacherStudentAssignment: assignments,
			School:                   schools,
			UserProfile:              profiles,
		},
		Integrations: &integrations.Integrations{Payments: pay},
	}
}

func subscriptionSchool(students int) *fakeSchools {
	return &fakeSchools{
		school:   &domain.School{ID: "s1", Billing: domain.SchoolBillingSubscription},
		students: students,
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
		schools     *fakeSchools
		profiles    *fakeProfiles
		wantRefused bool
	}{
		{
			name:        "there is room on the plan",
			assignments: &fakeAssignments{},
			schools:     subscriptionSchool(3),
			payments:    &fakePayments{entitlement: paidPlan(5)},
		},
		{
			name:        "the plan is full",
			assignments: &fakeAssignments{},
			schools:     subscriptionSchool(5),
			payments:    &fakePayments{entitlement: paidPlan(5)},
			wantRefused: true,
		},
		{
			// These paths are idempotent. Re-running one must not start failing
			// because the plan filled up in between; nothing is being added.
			name:        "a student the teacher already has is always allowed",
			assignments: &fakeAssignments{hasAccess: true},
			schools:     subscriptionSchool(99),
			payments:    &fakePayments{entitlement: paidPlan(5)},
		},
		{
			name:        "no subscription means the free allowance",
			assignments: &fakeAssignments{},
			schools:     subscriptionSchool(domain.FreePlan.MaxStudents),
			payments:    &fakePayments{entitlement: &payments.Entitlement{Active: false}},
			profiles:    withinTrial(),
			wantRefused: true,
		},
		{
			name:        "the free month still has room in it",
			assignments: &fakeAssignments{},
			schools:     subscriptionSchool(0),
			payments:    &fakePayments{entitlement: &payments.Entitlement{Active: false}},
			profiles:    withinTrial(),
		},
		{
			// The point of a trial. Once the free month is over the teacher
			// subscribes or adds nobody, even though they are under the one
			// student the free plan would otherwise allow.
			name:        "an expired trial allows nobody at all",
			assignments: &fakeAssignments{},
			schools:     subscriptionSchool(0),
			payments:    &fakePayments{entitlement: &payments.Entitlement{Active: false}},
			profiles:    pastTrial(),
			wantRefused: true,
		},
		{
			// Paying is what ends the trial's hold, so an expired one must not
			// follow a teacher who subscribed.
			name:        "a paid plan ignores the trial having run out",
			assignments: &fakeAssignments{},
			schools:     subscriptionSchool(3),
			payments:    &fakePayments{entitlement: paidPlan(5)},
			profiles:    pastTrial(),
		},
		{
			// The dangerous default. Treating an unreadable plan as the free
			// one would stop every paying teacher from working whenever the
			// payments service hiccups; one extra student costs far less.
			name:        "a payments outage does not block the teacher",
			assignments: &fakeAssignments{},
			schools:     subscriptionSchool(500),
			payments:    &fakePayments{err: errors.New("payments is down")},
		},
		{
			// Institutions are invoiced outside the product, so there is no
			// entitlement to read and nothing to cap.
			name:        "a school billed directly has no limit",
			assignments: &fakeAssignments{},
			schools: &fakeSchools{
				school:   &domain.School{ID: "s1", Billing: domain.SchoolBillingDirect},
				students: 5000,
			},
			payments: &fakePayments{entitlement: paidPlan(5)},
		},
		{
			// A teacher who only works at institutions owns no school, and an
			// institution's students are on nobody's subscription.
			name:        "a teacher with no school of their own is not capped",
			assignments: &fakeAssignments{},
			schools:     &fakeSchools{school: nil},
			payments:    &fakePayments{entitlement: paidPlan(5)},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := appWith(tc.assignments, tc.payments, tc.schools, tc.profiles)

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
