package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) listPages(ctx context.Context, notebookID string) ([]domain.NotebookPage, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, notebook_id, page_number, title, content_type, content_data, instructions, created_at
		FROM notebook_pages WHERE notebook_id = $1 ORDER BY page_number ASC
	`, notebookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []domain.NotebookPage
	for rows.Next() {
		var p domain.NotebookPage
		if err := rows.Scan(&p.ID, &p.NotebookID, &p.PageNumber, &p.Title, &p.ContentType, &p.ContentData, &p.Instructions, &p.CreatedAt); err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}
	return pages, nil
}
