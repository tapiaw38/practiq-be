package domain

import "time"

const DefaultTimezone = "America/Argentina/Buenos_Aires"

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

func StudentDay(at time.Time, loc *time.Location) time.Time {
	local := at.In(loc)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
}

func DaysBetween(from, to time.Time, loc *time.Location) int {
	fromDay := StudentDay(from, loc)
	toDay := StudentDay(to, loc)

	return int(toDay.Sub(fromDay).Round(24*time.Hour).Hours() / 24)
}

func EffectiveStreak(p StudentTopicProgress, loc *time.Location) int {
	if p.LastPracticedAt == nil {
		return 0
	}
	if DaysBetween(*p.LastPracticedAt, time.Now(), loc) <= 1 {
		return p.StreakDays
	}
	return 0
}

func CurrentStreak(progress []StudentTopicProgress, loc *time.Location) int {
	streak := 0
	for _, topic := range progress {
		if current := EffectiveStreak(topic, loc); current > streak {
			streak = current
		}
	}
	return streak
}
