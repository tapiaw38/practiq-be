package course

import "context"

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM courses WHERE id=$1`, id)
	return err
}
