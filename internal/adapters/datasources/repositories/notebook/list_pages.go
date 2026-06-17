package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) listPages(ctx context.Context, notebookID string) ([]domain.NotebookPage, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT np.id, np.notebook_id, np.page_number, np.title, np.content_type, np.content_data, np.instructions, np.created_at
		FROM notebook_pages np
		JOIN notebooks n ON n.id = np.notebook_id
		WHERE np.notebook_id = $1 AND n.deleted_at IS NULL
		ORDER BY np.page_number ASC
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
