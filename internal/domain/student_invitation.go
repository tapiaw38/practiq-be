package domain

import "time"

type StudentInvitation struct {
	ID        string
	Code      string
	TeacherID string

	SchoolID  string
	Uses      int
	ExpiresAt *time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

func (i StudentInvitation) IsUsable(now time.Time) bool {
	if i.RevokedAt != nil {
		return false
	}
	if i.ExpiresAt != nil && now.After(*i.ExpiresAt) {
		return false
	}

	return true
}
