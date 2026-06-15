package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) CreatePage(ctx context.Context, p domain.NotebookPage) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO notebook_pages (notebook_id, page_number, title, content_type, content_data, instructions)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`, p.NotebookID, p.PageNumber, p.Title, p.ContentType, p.ContentData, p.Instructions).Scan(&id)
	return id, err
}
