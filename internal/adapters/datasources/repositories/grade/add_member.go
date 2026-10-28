package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) AddMember(ctx context.Context, membership domain.GradeMembership) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO grade_memberships (grade_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (grade_id, user_id) DO NOTHING
	`, membership.GradeID, membership.UserID)
	return err
}
