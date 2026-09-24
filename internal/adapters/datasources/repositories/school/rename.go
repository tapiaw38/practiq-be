package school

import "context"

func (r *repository) Rename(ctx context.Context, schoolID, name string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE schools SET name = $2, updated_at = NOW() WHERE id = $1
	`, schoolID, name)
	return err
}
