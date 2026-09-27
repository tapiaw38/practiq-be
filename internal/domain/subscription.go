package domain

import "time"

var FreePlan = TeacherPlan{
	Name:        "Gratis",
	MaxStudents: 1,
	TrialDays:   30,
}

type (
	TeacherPlan struct {
		PlanID      int
		Name        string
		MaxStudents int

		TrialDays int
	}

	TeacherSubscription struct {
		Plan TeacherPlan

		Active bool

		StudentsUsed int

		RenewsAt *time.Time
	}
)

func TrialEndsAt(profileCreatedAt time.Time) time.Time {
	return profileCreatedAt.AddDate(0, 0, FreePlan.TrialDays)
}

func EffectiveFreePlan(profileCreatedAt, now time.Time) TeacherPlan {
	if now.After(TrialEndsAt(profileCreatedAt)) {
		expired := FreePlan
		expired.MaxStudents = 0
		return expired
	}
	return FreePlan
}

func PlanFromMetadata(planID int, name string, metadata map[string]any) TeacherPlan {
	plan := TeacherPlan{PlanID: planID, Name: name, MaxStudents: FreePlan.MaxStudents}
	if metadata == nil {
		return plan
	}

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

func (s TeacherSubscription) CanAddStudent() bool {
	return s.StudentsUsed < s.Plan.MaxStudents
}

func FreeSubscription(profileCreatedAt, now time.Time, studentsUsed int) TeacherSubscription {
	renewsAt := TrialEndsAt(profileCreatedAt)
	return TeacherSubscription{
		Plan:         EffectiveFreePlan(profileCreatedAt, now),
		Active:       false,
		StudentsUsed: studentsUsed,
		RenewsAt:     &renewsAt,
	}
}
