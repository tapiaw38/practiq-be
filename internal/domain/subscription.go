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

// TrialEndsAt is when a teacher's free month runs out. Counted from when their
// profile was created, which is the only moment that exists for every teacher.
func TrialEndsAt(profileCreatedAt time.Time) time.Time {
	return profileCreatedAt.AddDate(0, 0, FreePlan.TrialDays)
}

// EffectiveFreePlan is what a teacher who never paid may use right now.
//
// Once the free month is over it allows nothing. The trial is the whole free
// offer, not a discount on it: a teacher past it subscribes or keeps no
// students, which is the only reason to offer a month in the first place.
//
// Returning a zero allowance rather than a separate "expired" flag is what
// makes every consumer agree without being told: the limit check refuses the
// next student, the downgrade preview lists the ones over, and the screen
// shows 0 available. None of them need to know a trial exists.
func EffectiveFreePlan(profileCreatedAt, now time.Time) TeacherPlan {
	if now.After(TrialEndsAt(profileCreatedAt)) {
		expired := FreePlan
		expired.MaxStudents = 0
		return expired
	}
	return FreePlan
}

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

// FreeSubscription is the answer for a teacher who has never paid.
//
// RenewsAt is the day the free month ends whether or not it already has, so a
// teacher whose trial ran out is told when rather than left with a blank date.
func FreeSubscription(profileCreatedAt, now time.Time, studentsUsed int) TeacherSubscription {
	renewsAt := TrialEndsAt(profileCreatedAt)
	return TeacherSubscription{
		Plan:         EffectiveFreePlan(profileCreatedAt, now),
		Active:       false,
		StudentsUsed: studentsUsed,
		RenewsAt:     &renewsAt,
	}
}
