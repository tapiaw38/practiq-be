package studentinvitation

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Create(ctx context.Context, invitation domain.StudentInvitation) (*domain.StudentInvitation, error) {
	query := `
		INSERT INTO student_invitations (code, teacher_id, expires_at, school_id)
		VALUES ($1, $2, $3, NULLIF($4, '')::uuid)
		RETURNING id, code, teacher_id, uses, expires_at, revoked_at, created_at
	`

	row := r.db.QueryRowContext(ctx, query, invitation.Code, invitation.TeacherID, invitation.ExpiresAt, tenant.SchoolID(ctx))

	return scanInvitation(row)
}
