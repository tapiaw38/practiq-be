package domain

import "time"

// GraceDays is how long a school keeps the plan it stopped paying for.
//
// Somebody whose card expires on a Friday should not find their class locked
// out on Saturday morning. Five days is long enough to notice an email and fix
// a card, and short enough that it is not a free month.
const GraceDays = 5

// GraceEndsAt is when a lapsed plan stops being honoured.
func GraceEndsAt(paidUntil time.Time) time.Time {
	return paidUntil.AddDate(0, 0, GraceDays)
}

// InGrace reports whether a plan that stopped being paid should still be
// honoured in full.
//
// False before the paid period ends too: that is not grace, it is simply paid,
// and the caller has an entitlement for it. Only the window after counts.
func InGrace(paidUntil time.Time, now time.Time) bool {
	if paidUntil.IsZero() {
		return false
	}
	return !now.Before(paidUntil) && now.Before(GraceEndsAt(paidUntil))
}
