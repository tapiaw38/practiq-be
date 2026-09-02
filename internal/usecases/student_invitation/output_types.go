package studentinvitation

import (
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/invitecode"
)

type InvitationData struct {
	ID string `json:"id"`
	// Code viaja sin guion y FormattedCode con él: el alumno escribe
	// cualquiera de los dos y la pantalla muestra el legible.
	Code          string     `json:"code"`
	FormattedCode string     `json:"formatted_code"`
	Uses          int        `json:"uses"`
	ExpiresAt     *time.Time `json:"expires_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

func toInvitationData(invitation domain.StudentInvitation) InvitationData {
	return InvitationData{
		ID:            invitation.ID,
		Code:          invitation.Code,
		FormattedCode: invitecode.Format(invitation.Code),
		Uses:          invitation.Uses,
		ExpiresAt:     invitation.ExpiresAt,
		CreatedAt:     invitation.CreatedAt,
	}
}
