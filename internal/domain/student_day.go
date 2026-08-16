package domain

import "time"

// DefaultTimezone is where the product started. It is the fallback when a
// student has no timezone recorded, which is every student created before the
// field existed.
const DefaultTimezone = "America/Argentina/Buenos_Aires"

// StudentLocation resolves the zone a student's day is measured in.
//
// A streak counts calendar days, so "which day is it" has to be answered in
// the student's own time. Measuring in UTC made the day flip at 21:00 in
// Argentina: two sessions the same evening counted as two days, and two
// consecutive evenings counted as one. Evening is when homework happens, so
// that was the common case, not an edge one.
func StudentLocation(timezone string) *time.Location {
	if timezone != "" {
		if loc, err := time.LoadLocation(timezone); err == nil {
			return loc
		}
	}
	if loc, err := time.LoadLocation(DefaultTimezone); err == nil {
		return loc
	}
	return time.UTC
}

// StudentDay truncates an instant to midnight in the student's zone, so two
// instants can be compared as calendar days.
func StudentDay(at time.Time, loc *time.Location) time.Time {
	local := at.In(loc)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
}

// DaysBetween counts calendar days from one instant to another in the
// student's zone. Negative when `to` precedes `from`.
func DaysBetween(from, to time.Time, loc *time.Location) int {
	fromDay := StudentDay(from, loc)
	toDay := StudentDay(to, loc)
	// Hours/24 rather than subtracting dates: it stays correct across DST
	// shifts, where a local day is 23 or 25 hours long.
	return int(toDay.Sub(fromDay).Round(24*time.Hour).Hours() / 24)
}

// EffectiveStreak hides a streak the student already broke.
//
// The stored counter is only recomputed when that topic is practised again, so
// on its own it keeps whatever value a student had when they stopped. Anything
// that shows a streak has to go through here, or a report ends up claiming 12
// days for someone who has not practised in a month.
func EffectiveStreak(p StudentTopicProgress, loc *time.Location) int {
	if p.LastPracticedAt == nil {
		return 0
	}
	if DaysBetween(*p.LastPracticedAt, time.Now(), loc) <= 1 {
		return p.StreakDays
	}
	return 0
}
