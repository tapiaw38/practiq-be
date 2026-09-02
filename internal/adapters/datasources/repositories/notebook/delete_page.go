package notebook

import "context"

func (r *repository) DeletePage(ctx context.Context, pageID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM notebook_pages WHERE id = $1`, pageID)
	return err
}
