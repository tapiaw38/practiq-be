package studentcoursexp

import (
	"context"
	"database/sql"
	"errors"
)

// Award records the event before increasing the balance. The unique event key
// makes retries harmless: a request can be repeated without farming XP.
func (r *repository) Award(ctx context.Context, input AwardInput) (result AwardResult, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var points int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO student_course_xp_events
			(student_id, course_id, event_key, event_type, points, practice_sheet_id, exercise_id)
		VALUES ($1, $2::uuid, $3, $4, $5, NULLIF($6, '')::uuid, NULLIF($7, '')::uuid)
		ON CONFLICT (student_id, course_id, event_key) DO NOTHING
		RETURNING points`,
		input.StudentID, input.CourseID, input.EventKey, input.EventType, input.Points,
		input.PracticeSheetID, input.ExerciseID,
	).Scan(&points)
	if errors.Is(err, sql.ErrNoRows) {
		// The event key already existed, so this is a retry: report the
		// balance without awarding again. A missing balance row is a zero
		// balance, not a failure.
		err = tx.QueryRowContext(ctx, `
			SELECT total_xp FROM student_course_xp
			WHERE student_id = $1 AND course_id = $2::uuid`, input.StudentID, input.CourseID,
		).Scan(&result.TotalXP)
		if errors.Is(err, sql.ErrNoRows) {
			result.TotalXP = 0
			err = nil
		}
		if err != nil {
			return result, err
		}
		if err = tx.Commit(); err != nil {
			return result, err
		}
		return result, nil
	}
	if err != nil {
		return result, err
	}

	err = tx.QueryRowContext(ctx, `
		INSERT INTO student_course_xp (student_id, course_id, total_xp, created_at, updated_at)
		VALUES ($1, $2::uuid, $3, now(), now())
		ON CONFLICT (student_id, course_id) DO UPDATE
		SET total_xp = student_course_xp.total_xp + EXCLUDED.total_xp,
			updated_at = now()
		RETURNING total_xp`, input.StudentID, input.CourseID, points,
	).Scan(&result.TotalXP)
	if err != nil {
		return result, err
	}
	result.Awarded = true
	if err = tx.Commit(); err != nil {
		return result, err
	}
	return result, nil
}
