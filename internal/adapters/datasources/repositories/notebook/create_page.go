package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) CreatePage(ctx context.Context, p domain.NotebookPage) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO notebook_pages (notebook_id, page_number, title, content_type, content_data, statement_text, statement_verified, instructions, school_id)
		SELECT $1, $2, $3, $4, $5, $6, $7, $8, NULLIF($9,'')::uuid
		WHERE EXISTS (SELECT 1 FROM notebooks WHERE id = $1 AND ($9 = '' OR school_id = NULLIF($9,'')::uuid))
		RETURNING id
	`, p.NotebookID, p.PageNumber, p.Title, p.ContentType, p.ContentData, p.StatementText, p.StatementVerified, p.Instructions, tenant.SchoolID(ctx)).Scan(&id)
	return id, err
}
