package grade

import (
	"context"

	"github.com/lib/pq"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

// List returns the grades of the given schools. A nil slice means no narrowing,
// which is a platform superadmin — not "no schools", which returns nothing.
func (r *repository) List(ctx context.Context, schoolIDs []string) ([]domain.Grade, error) {
	query := `
		SELECT id, name, COALESCE(description, ''), visual_theme, created_by, created_at
		FROM grades
	`
	args := []any{}
	if schoolIDs != nil {
		query += ` WHERE school_id = ANY($1)`
		args = append(args, pq.Array(schoolIDs))
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grades := []domain.Grade{}
	for rows.Next() {
		var grade domain.Grade
		if err := rows.Scan(&grade.ID, &grade.Name, &grade.Description, &grade.VisualTheme, &grade.CreatedBy, &grade.CreatedAt); err != nil {
			return nil, err
		}
		grades = append(grades, grade)
	}
	return grades, nil
}
