package notebook

import (
	"context"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) DeletePage(ctx context.Context, pageID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM notebook_pages WHERE id = $1 AND ($2 = '' OR school_id = NULLIF($2,'')::uuid)`, pageID, tenant.SchoolID(ctx))
	return err
}
