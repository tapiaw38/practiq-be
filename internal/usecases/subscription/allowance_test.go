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

	school *domain.School

	byID     map[string]*domain.School
	students int
	memberOf []domain.SchoolMember

	studentsByActivity []string
}

func (f *fakeSchools) ListStudentsByActivity(context.Context, string) ([]string, error) {
	return f.studentsByActivity, nil
}

func (f *fakeSchools) GetPersonal(context.Context, string) (*domain.School, error) {
	return f.school, nil
}

func (f *fakeSchools) Get(_ context.Context, id string) (*domain.School, error) {
	return f.byID[id], nil
}

func (f *fakeSchools) CountStudents(context.Context, string) (int, error) {
	return f.students, nil
}

func (f *fakeSchools) ListForUser(context.Context, string) ([]domain.SchoolMember, error) {
	return f.memberOf, nil
}

type fakePayments struct {
	payments.Client
	entitlement   *payments.Entitlement
	subscriptions []payments.Subscription
	plans         []payments.Plan
	err           error
}

func (f *fakePayments) GetEntitlement(context.Context, string) (*payments.Entitlement, error) {
	return f.entitlement, f.err
}

func (f *fakePayments) ListSubscriptions(context.Context, string) ([]payments.Subscription, error) {
	return f.subscriptions, f.err
}

func (f *fakePayments) ListPlans(context.Context) ([]payments.Plan, error) {
	return f.plans, f.err
}

type fakeProfiles struct {
	userprofile.Repository
	createdAt time.Time
}

func (f *fakeProfiles) Get(context.Context, string) (*domain.UserProfile, error) {
	return &domain.UserProfile{CreatedAt: f.createdAt}, nil
}

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
		school:   &domain.School{ID: "s1", Kind: domain.SchoolKindPersonal, Billing: domain.SchoolBillingSubscription, CreatedBy: "teacher-1"},
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

		schoolID    string
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

			name:        "a student already in the school is always allowed",
			assignments: &fakeAssignments{},
			schools: func() *fakeSchools {
				s := subscriptionSchool(99)
				s.memberOf = []domain.SchoolMember{{SchoolID: "s1", Role: domain.SchoolRoleStudent, Active: true}}
				return s
			}(),
			payments: &fakePayments{entitlement: paidPlan(5)},
		},
		{

			name:        "a deactivated student is not already in, so the plan decides",
			assignments: &fakeAssignments{hasAccess: true},
			schools: func() *fakeSchools {
				s := subscriptionSchool(5)
				s.memberOf = []domain.SchoolMember{{SchoolID: "s1", Role: domain.SchoolRoleStudent, Active: false}}
				return s
			}(),
			payments:    &fakePayments{entitlement: paidPlan(5)},
			wantRefused: true,
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

			name:        "an expired trial allows nobody at all",
			assignments: &fakeAssignments{},
			schools:     subscriptionSchool(0),
			payments:    &fakePayments{entitlement: &payments.Entitlement{Active: false}},
			profiles:    pastTrial(),
			wantRefused: true,
		},
		{

			name:        "a paid plan ignores the trial having run out",
			assignments: &fakeAssignments{},
			schools:     subscriptionSchool(3),
			payments:    &fakePayments{entitlement: paidPlan(5)},
			profiles:    pastTrial(),
		},
		{

			name:        "a payments outage does not block the teacher",
			assignments: &fakeAssignments{},
			schools:     subscriptionSchool(500),
			payments:    &fakePayments{err: errors.New("payments is down")},
		},
		{

			name:        "a school billed directly has no limit",
			assignments: &fakeAssignments{},
			schools: &fakeSchools{
				school:   &domain.School{ID: "s1", Billing: domain.SchoolBillingDirect},
				students: 5000,
			},
			payments: &fakePayments{entitlement: paidPlan(5)},
		},
		{

			name:        "a teacher with no school of their own is not capped",
			assignments: &fakeAssignments{},
			schools:     &fakeSchools{school: nil},
			payments:    &fakePayments{entitlement: paidPlan(5)},
		},
		{

			name:        "a full personal plan does not refuse an institution's student",
			assignments: &fakeAssignments{},
			schools: &fakeSchools{
				school:   &domain.School{ID: "personal", Kind: domain.SchoolKindPersonal, Billing: domain.SchoolBillingSubscription, CreatedBy: "teacher-1"},
				byID:     map[string]*domain.School{"inst": {ID: "inst", Kind: domain.SchoolKindInstitution, Billing: domain.SchoolBillingDirect}},
				students: 5,
			},
			payments: &fakePayments{entitlement: paidPlan(5)},
			schoolID: "inst",
		},
		{

			name:        "an institution is uncapped even when billed by subscription",
			assignments: &fakeAssignments{},
			schools: &fakeSchools{
				school:   &domain.School{ID: "personal", Kind: domain.SchoolKindPersonal, Billing: domain.SchoolBillingSubscription, CreatedBy: "teacher-1"},
				byID:     map[string]*domain.School{"inst": {ID: "inst", Kind: domain.SchoolKindInstitution, Billing: domain.SchoolBillingSubscription}},
				students: 5000,
			},
			payments: &fakePayments{entitlement: paidPlan(5)},
			schoolID: "inst",
		},
		{

			name:        "naming the personal school still applies its plan",
			assignments: &fakeAssignments{},
			schools: &fakeSchools{
				school:   &domain.School{ID: "personal", Kind: domain.SchoolKindPersonal, Billing: domain.SchoolBillingSubscription, CreatedBy: "teacher-1"},
				byID:     map[string]*domain.School{"personal": {ID: "personal", Kind: domain.SchoolKindPersonal, Billing: domain.SchoolBillingSubscription, CreatedBy: "teacher-1"}},
				students: 5,
			},
			payments:    &fakePayments{entitlement: paidPlan(5)},
			schoolID:    "personal",
			wantRefused: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := appWith(tc.assignments, tc.payments, tc.schools, tc.profiles)

			appErr := EnsureCanAddStudent(context.Background(), app, tc.schoolID, "teacher-1", "student-1")

			if tc.wantRefused && appErr == nil {
				t.Fatal("expected the link to be refused")
			}
			if !tc.wantRefused && appErr != nil {
				t.Fatalf("expected the link to be allowed, got %v", appErr)
			}
		})
	}
}

func TestGraceKeepsTheLapsedPlan(t *testing.T) {
	lapsedYesterday := payments.Timestamp{Time: time.Now().UTC().AddDate(0, 0, -1)}
	lapsedLongAgo := payments.Timestamp{Time: time.Now().UTC().AddDate(0, 0, -domain.GraceDays-1)}

	cases := []struct {
		name            string
		periodEnd       *payments.Timestamp
		wantMaxStudents int
	}{
		{"lapsed yesterday keeps the plan", &lapsedYesterday, 15},

		{"grace spent falls to the expired trial", &lapsedLongAgo, 0},
		{"no recorded period is not grace", nil, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := appWith(nil, &fakePayments{

				entitlement: &payments.Entitlement{Active: false},
				subscriptions: []payments.Subscription{
					{ID: 10, PlanID: 2, Status: "cancelled", CurrentPeriodEnd: tc.periodEnd},
				},
				plans: []payments.Plan{
					{ID: 2, Metadata: map[string]any{"name": "Crecimiento", "max_students": float64(15)}},
				},
			}, subscriptionSchool(0), pastTrial())

			scope, appErr := scopeFor(t.Context(), app, "", "teacher-1")
			if appErr != nil {
				t.Fatalf("scopeFor: %v", appErr)
			}
			if scope.Plan.MaxStudents != tc.wantMaxStudents {
				t.Fatalf("max students = %d, want %d", scope.Plan.MaxStudents, tc.wantMaxStudents)
			}

			if scope.Active {
				t.Fatal("a lapsed plan must not read as active")
			}
		})
	}
}

func TestStudentsCanWork(t *testing.T) {
	inGrace := payments.Timestamp{Time: time.Now().UTC().AddDate(0, 0, -1)}

	cases := []struct {
		name     string
		pay      *fakePayments
		profiles *fakeProfiles
		want     bool
	}{
		{
			name:     "paying",
			pay:      &fakePayments{entitlement: &payments.Entitlement{Active: true, Metadata: map[string]any{"max_students": float64(5)}}},
			profiles: pastTrial(),
			want:     true,
		},
		{
			name: "lapsed but inside the grace window",
			pay: &fakePayments{
				entitlement:   &payments.Entitlement{Active: false},
				subscriptions: []payments.Subscription{{ID: 1, PlanID: 2, Status: "cancelled", CurrentPeriodEnd: &inGrace}},
				plans:         []payments.Plan{{ID: 2, Metadata: map[string]any{"max_students": float64(15)}}},
			},
			profiles: pastTrial(),
			want:     true,
		},
		{
			name:     "trial still running",
			pay:      &fakePayments{entitlement: &payments.Entitlement{Active: false}},
			profiles: withinTrial(),
			want:     true,
		},
		{
			name:     "trial spent and nobody paying",
			pay:      &fakePayments{entitlement: &payments.Entitlement{Active: false}},
			profiles: pastTrial(),
			want:     false,
		},
		{
			name:     "a payments outage is not a failure to pay",
			pay:      &fakePayments{err: errors.New("payments down")},
			profiles: pastTrial(),
			want:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := appWith(nil, tc.pay, subscriptionSchool(0), tc.profiles)
			if got := StudentsCanWork(t.Context(), app, "", "teacher-1"); got != tc.want {
				t.Fatalf("StudentsCanWork = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStudentMayWorkAppliesTheCapBeforeTheTeacherChooses(t *testing.T) {
	paidForFive := &fakePayments{
		entitlement: &payments.Entitlement{Active: true, Metadata: map[string]any{"max_students": float64(2)}},
	}

	schools := subscriptionSchool(4)
	schools.studentsByActivity = []string{"dormant-1", "dormant-2", "working-1", "working-2"}

	app := appWith(nil, paidForFive, schools, pastTrial())

	for _, tc := range []struct {
		student string
		want    bool
	}{
		{"working-1", true},
		{"working-2", true},
		{"dormant-1", false},
		{"dormant-2", false},
	} {
		t.Run(tc.student, func(t *testing.T) {
			if got := StudentMayWork(t.Context(), app, "", "teacher-1", tc.student); got != tc.want {
				t.Fatalf("StudentMayWork(%s) = %v, want %v", tc.student, got, tc.want)
			}
		})
	}
}
