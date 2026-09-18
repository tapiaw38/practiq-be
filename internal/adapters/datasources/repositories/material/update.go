package material

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, m domain.Material) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE materials SET title = $1, extracted_text = $2, file_url = $3, status = $4 WHERE id = $5
	`, m.Title, m.ExtractedText, m.FileURL, m.Status, id)
	return err
}
