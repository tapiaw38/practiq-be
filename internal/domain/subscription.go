package domain

import "time"

// FreePlan is what a teacher has before paying for anything: one student, for
// a month.
//
// It is not a plan in the payments service, and deliberately so. Creating a
// gateway subscription that charges nothing, only to not charge it, buys
// nothing and puts a teacher's first month at the mercy of a payment provider
// being reachable. Payments gets involved when money does.
var FreePlan = TeacherPlan{
	Name:        "Gratis",
	MaxStudents: 1,
	TrialDays:   30,
}

type (
	// TeacherPlan is what a plan allows, in Practiq's terms.
	//
	// The payments service stores this as opaque metadata on the plan and never
	// reads it, so the meaning of these fields is decided here and can change
	// without touching that service.
	TeacherPlan struct {
		// PlanID is empty on the free plan, which payments does not know about.
		PlanID      int
		Name        string
		MaxStudents int
		// TrialDays only applies to the free plan.
		TrialDays int
	}

	// TeacherSubscription is what the teacher's subscription tab shows, and
	// what the student-limit check reads.
	TeacherSubscription struct {
		Plan TeacherPlan
		// Active is false while the teacher is on the free plan, whether or not
		// the trial has run out.
		Active bool
		// StudentsUsed counts by the same rule the rest of the product uses to
		// decide who a teacher's students are.
		StudentsUsed int
		// RenewsAt is when the paid period ends, or when the free month does.
		RenewsAt *time.Time
	}
)

// PlanFromMetadata reads a payments plan's metadata into Practiq's terms.
//
// Anything missing falls back to the free plan's allowance rather than to
// unlimited: a plan saved without a limit should sell nothing, not everything.
func PlanFromMetadata(planID int, name string, metadata map[string]any) TeacherPlan {
	plan := TeacherPlan{PlanID: planID, Name: name, MaxStudents: FreePlan.MaxStudents}
	if metadata == nil {
		return plan
	}
	// JSON numbers arrive as float64 through an untyped map.
	if raw, ok := metadata["max_students"]; ok {
		switch value := raw.(type) {
		case float64:
			plan.MaxStudents = int(value)
		case int:
			plan.MaxStudents = value
		}
	}
	return plan
}

// CanAddStudent reports whether one more student fits.
func (s TeacherSubscription) CanAddStudent() bool {
	return s.StudentsUsed < s.Plan.MaxStudents
}

// FreeSubscription is the answer for a teacher who has never paid. The free
// month is counted from when their profile was created, which is the only
// moment that exists for every teacher.
func FreeSubscription(profileCreatedAt time.Time, studentsUsed int) TeacherSubscription {
	renewsAt := profileCreatedAt.AddDate(0, 0, FreePlan.TrialDays)
	return TeacherSubscription{
		Plan:         FreePlan,
		Active:       false,
		StudentsUsed: studentsUsed,
		RenewsAt:     &renewsAt,
	}
}
