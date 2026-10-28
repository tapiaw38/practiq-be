package notebook

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) GetPage(ctx context.Context, pageID string) (*domain.NotebookPage, error) {
	var p domain.NotebookPage
	err := r.db.QueryRowContext(ctx, `
		SELECT np.id, np.notebook_id, np.page_number, np.title, np.content_type, np.content_data, np.statement_text, np.statement_verified, np.instructions, np.created_at
		FROM notebook_pages np
		JOIN notebooks n ON n.id = np.notebook_id
		WHERE np.id = $1 AND n.deleted_at IS NULL
	`, pageID).Scan(&p.ID, &p.NotebookID, &p.PageNumber, &p.Title, &p.ContentType, &p.ContentData, &p.StatementText, &p.StatementVerified, &p.Instructions, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}
