package material

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Update(ctx context.Context, id string, m domain.Material) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE materials SET title = $1, extracted_text = $2, file_url = $3, status = $4 WHERE id = $5 AND ($6 = '' OR school_id = NULLIF($6, '')::uuid)
	`, m.Title, m.ExtractedText, m.FileURL, m.Status, id, tenant.SchoolID(ctx))
	return err
}
