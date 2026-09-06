package practicesheet

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.PracticeSheet, error) {
	query := `
		SELECT ps.id, ps.course_id, COALESCE(ps.topic_id::text,''), COALESCE(ps.strategy_id::text,''), ps.title, ps.level, ps.sheet_type, ps.test_style, ps.scheduled_at, ps.available_until, ps.created_by, ps.created_at
		FROM practice_sheets ps
		JOIN courses c ON c.id = ps.course_id
		WHERE ps.id = $1 AND c.deleted_at IS NULL AND ($2 = '' OR c.school_id = NULLIF($2, '')::uuid)
	`
	row := r.db.QueryRowContext(ctx, query, id, tenant.SchoolID(ctx))
	var ps domain.PracticeSheet
	err := row.Scan(&ps.ID, &ps.CourseID, &ps.TopicID, &ps.StrategyID, &ps.Title, &ps.Level, &ps.SheetType, &ps.TestStyle, &ps.ScheduledAt, &ps.AvailableUntil, &ps.CreatedBy, &ps.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Load exercises
	exQuery := `
		SELECT pse.id, pse.practice_sheet_id, pse.order_index,
		       e.id, e.topic_id, COALESCE(e.material_id::text,''), e.type, e.question,
		       COALESCE(e.correct_answer,''), COALESCE(e.explanation,''), e.difficulty, e.metadata::text, e.created_at
		FROM practice_sheet_exercises pse
		JOIN exercises e ON e.id = pse.exercise_id
		JOIN practice_sheets ps ON ps.id = pse.practice_sheet_id
		WHERE pse.practice_sheet_id = $1 AND ($2 = '' OR ps.school_id = NULLIF($2, '')::uuid)
		ORDER BY pse.order_index ASC
	`
	rows, err := r.db.QueryContext(ctx, exQuery, id, tenant.SchoolID(ctx))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pse domain.PracticeSheetExercise
		if err := rows.Scan(
			&pse.ID, &pse.PracticeSheetID, &pse.OrderIndex,
			&pse.Exercise.ID, &pse.Exercise.TopicID, &pse.Exercise.MaterialID, &pse.Exercise.Type,
			&pse.Exercise.Question, &pse.Exercise.CorrectAnswer, &pse.Exercise.Explanation,
			&pse.Exercise.Difficulty, &pse.Exercise.Metadata, &pse.Exercise.CreatedAt,
		); err != nil {
			return nil, err
		}
		ps.Exercises = append(ps.Exercises, pse)
	}
	return &ps, nil
}
