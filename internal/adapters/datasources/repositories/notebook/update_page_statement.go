package notebook

import (
	"context"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) UpdatePageStatement(ctx context.Context, pageID, statementText string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notebook_pages
		SET statement_text = $1
		WHERE id = $2 AND ($3 = '' OR school_id = NULLIF($3,'')::uuid)
	`, statementText, pageID, tenant.SchoolID(ctx))
	return err
}
