package domain

import "time"

type StudentInvitation struct {
	ID        string
	Code      string
	TeacherID string
	Uses      int
	ExpiresAt *time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// IsUsable responde si el código todavía sirve para vincularse. La revocación
// gana sobre el vencimiento: un código que el docente dio de baja no vuelve a
// servir aunque se le extienda la fecha.
func (i StudentInvitation) IsUsable(now time.Time) bool {
	if i.RevokedAt != nil {
		return false
	}
	if i.ExpiresAt != nil && now.After(*i.ExpiresAt) {
		return false
	}

	return true
}
