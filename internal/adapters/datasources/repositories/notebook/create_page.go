package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) CreatePage(ctx context.Context, p domain.NotebookPage) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO notebook_pages (notebook_id, page_number, title, content_type, content_data, statement_text, statement_verified, instructions)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id
	`, p.NotebookID, p.PageNumber, p.Title, p.ContentType, p.ContentData, p.StatementText, p.StatementVerified, p.Instructions).Scan(&id)
	return id, err
}
