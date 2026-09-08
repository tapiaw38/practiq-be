package domain

import (
	"testing"
	"time"
)

func TestPlanFromMetadataReadsTheLimitTheProductPublished(t *testing.T) {
	cases := []struct {
		name     string
		metadata map[string]any
		want     int
	}{
		{
			// The payments service stores metadata as JSON, so numbers arrive
			// as float64 through an untyped map.
			name:     "a limit sent as JSON is read",
			metadata: map[string]any{"max_students": float64(5)},
			want:     5,
		},
		{
			// Falling back to the free allowance rather than to unlimited: a
			// plan saved without a limit should sell nothing, not everything.
			name:     "a plan with no limit grants the free allowance",
			metadata: map[string]any{"name": "Equipo"},
			want:     FreePlan.MaxStudents,
		},
		{
			name:     "metadata that never arrived grants the free allowance",
			metadata: nil,
			want:     FreePlan.MaxStudents,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := PlanFromMetadata(7, "Equipo", tc.metadata)
			if got.MaxStudents != tc.want {
				t.Fatalf("max students = %d, want %d", got.MaxStudents, tc.want)
			}
		})
	}
}

func TestCanAddStudentStopsAtTheLimit(t *testing.T) {
	plan := TeacherPlan{MaxStudents: 5}

	cases := []struct {
		used int
		want bool
	}{
		{used: 0, want: true},
		{used: 4, want: true},
		// The fifth student fills the plan; a sixth does not fit.
		{used: 5, want: false},
		// Already over, which happens after a downgrade.
		{used: 9, want: false},
	}

	for _, tc := range cases {
		got := TeacherSubscription{Plan: plan, StudentsUsed: tc.used}.CanAddStudent()
		if got != tc.want {
			t.Fatalf("with %d of %d used: got %v, want %v", tc.used, plan.MaxStudents, got, tc.want)
		}
	}
}

func TestFreeSubscriptionRunsAMonthFromTheProfile(t *testing.T) {
	created := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)

	got := FreeSubscription(created, 1)

	if got.Active {
		t.Fatal("the free plan is not an active subscription")
	}
	if got.Plan.MaxStudents != 1 {
		t.Fatalf("free plan allows %d students, want 1", got.Plan.MaxStudents)
	}
	if got.CanAddStudent() {
		t.Fatal("one student already fills the free plan")
	}
	want := created.AddDate(0, 0, 30)
	if got.RenewsAt == nil || !got.RenewsAt.Equal(want) {
		t.Fatalf("renews at %v, want %v", got.RenewsAt, want)
	}
}
