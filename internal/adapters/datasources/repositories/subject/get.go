package subject

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.Subject, error) {
	var subject domain.Subject
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, COALESCE(description, ''), created_by, created_at
		FROM subjects
		WHERE id = $1
	`, id).Scan(&subject.ID, &subject.Name, &subject.Description, &subject.CreatedBy, &subject.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &subject, nil
}
