package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) UpdatePage(ctx context.Context, p domain.NotebookPage) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notebook_pages
		SET title = $1, content_type = $2, content_data = $3, statement_text = $4, statement_verified = $5, instructions = $6
		WHERE id = $7
	`, p.Title, p.ContentType, p.ContentData, p.StatementText, p.StatementVerified, p.Instructions, p.ID)
	return err
}
