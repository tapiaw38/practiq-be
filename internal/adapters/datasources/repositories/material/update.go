package material

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, m domain.Material) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE materials SET title = $1, extracted_text = $2 WHERE id = $3
	`, m.Title, m.ExtractedText, id)
	return err
}
