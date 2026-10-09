package practicesheet

import "context"

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM practice_sheets WHERE id = $1`, id)
	return err
}
