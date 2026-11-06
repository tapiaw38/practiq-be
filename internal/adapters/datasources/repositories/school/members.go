package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) RemoveMember(ctx context.Context, schoolID, userID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM school_members WHERE school_id = $1 AND user_id = $2`, schoolID, userID)
	return err
}

func (r *repository) CountActiveAdmins(ctx context.Context, schoolID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM school_members WHERE school_id = $1 AND role = 'admin' AND active`, schoolID).Scan(&count)
	return count, err
}

func (r *repository) ListMembers(ctx context.Context, schoolID string) ([]domain.SchoolMember, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT school_id, user_id, role, active
		FROM school_members WHERE school_id = $1
		ORDER BY role, created_at
	`, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []domain.SchoolMember{}
	for rows.Next() {
		var m domain.SchoolMember
		if err := rows.Scan(&m.SchoolID, &m.UserID, &m.Role, &m.Active); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *repository) CountStudents(ctx context.Context, schoolID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM school_members WHERE school_id = $1 AND role = 'student' AND active`, schoolID).Scan(&count)
	return count, err
}

func (r *repository) SetMemberActive(ctx context.Context, schoolID, userID string, active bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE school_members SET active = $3 WHERE school_id = $1 AND user_id = $2`, schoolID, userID, active)
	return err
}
