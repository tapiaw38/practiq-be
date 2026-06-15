package notebook

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) GetPage(ctx context.Context, pageID string) (*domain.NotebookPage, error) {
	var p domain.NotebookPage
	err := r.db.QueryRowContext(ctx, `
		SELECT id, notebook_id, page_number, title, content_type, content_data, instructions, created_at
		FROM notebook_pages WHERE id = $1
	`, pageID).Scan(&p.ID, &p.NotebookID, &p.PageNumber, &p.Title, &p.ContentType, &p.ContentData, &p.Instructions, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}
