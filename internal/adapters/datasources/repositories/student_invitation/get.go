package studentinvitation

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

const selectColumns = `id, code, teacher_id, uses, expires_at, revoked_at, created_at`

func (r *repository) GetByCode(ctx context.Context, code string) (*domain.StudentInvitation, error) {
	query := `SELECT ` + selectColumns + ` FROM student_invitations WHERE code = $1 AND ($2 = '' OR school_id = NULLIF($2, '')::uuid)`

	return scanInvitation(r.db.QueryRowContext(ctx, query, code, tenant.SchoolID(ctx)))
}

func (r *repository) GetActiveByTeacher(ctx context.Context, teacherID string) (*domain.StudentInvitation, error) {
	query := `
		SELECT ` + selectColumns + `
		FROM student_invitations
		WHERE teacher_id = $1 AND revoked_at IS NULL
		  AND ($2 = '' OR school_id = NULLIF($2, '')::uuid)
	`

	return scanInvitation(r.db.QueryRowContext(ctx, query, teacherID, tenant.SchoolID(ctx)))
}

type scanner interface {
	Scan(dest ...any) error
}

func scanInvitation(row scanner) (*domain.StudentInvitation, error) {
	var invitation domain.StudentInvitation
	err := row.Scan(
		&invitation.ID,
		&invitation.Code,
		&invitation.TeacherID,
		&invitation.Uses,
		&invitation.ExpiresAt,
		&invitation.RevokedAt,
		&invitation.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &invitation, nil
}
