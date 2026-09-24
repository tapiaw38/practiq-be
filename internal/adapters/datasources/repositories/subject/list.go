package subject

import (
	"context"

	"github.com/lib/pq"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

// List returns the subjects of the given schools. A nil slice means no
// narrowing, which is a platform superadmin — not "no schools", which returns
// nothing.
func (r *repository) List(ctx context.Context, schoolIDs []string) ([]domain.Subject, error) {
	query := `
		SELECT id, name, COALESCE(description, ''), created_by, created_at
		FROM subjects
	`
	args := []any{}
	if schoolIDs != nil {
		query += ` WHERE school_id = ANY($1)`
		args = append(args, pq.Array(schoolIDs))
	}
	query += ` ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subjects := []domain.Subject{}
	for rows.Next() {
		var subject domain.Subject
		if err := rows.Scan(&subject.ID, &subject.Name, &subject.Description, &subject.CreatedBy, &subject.CreatedAt); err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	return subjects, nil
}
