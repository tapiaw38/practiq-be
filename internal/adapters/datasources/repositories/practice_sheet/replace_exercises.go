package practicesheet

import "context"

func (r *repository) ReplaceExercises(ctx context.Context, sheetID string, exerciseIDs []string) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM practice_sheet_exercises WHERE practice_sheet_id = $1`, sheetID); err != nil {
		return err
	}
	for i, eid := range exerciseIDs {
		if err := r.AddExercise(ctx, sheetID, eid, i); err != nil {
			return err
		}
	}
	return nil
}
