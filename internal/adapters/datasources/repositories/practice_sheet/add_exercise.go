package practicesheet

import "context"

func (r *repository) AddExercise(ctx context.Context, sheetID, exerciseID string, orderIndex int) error {
	query := `INSERT INTO practice_sheet_exercises (practice_sheet_id, exercise_id, order_index) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, sheetID, exerciseID, orderIndex)
	return err
}
