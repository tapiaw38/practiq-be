package domain

import "time"

const GraceDays = 5

func GraceEndsAt(paidUntil time.Time) time.Time {
	return paidUntil.AddDate(0, 0, GraceDays)
}

func InGrace(paidUntil time.Time, now time.Time) bool {
	if paidUntil.IsZero() {
		return false
	}
	return !now.Before(paidUntil) && now.Before(GraceEndsAt(paidUntil))
}
