package practicesheet

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) List(ctx context.Context, filter ListFilter) ([]domain.PracticeSheet, error) {
	query := `
		SELECT id, course_id, COALESCE(topic_id::text,''), COALESCE(strategy_id::text,''), title, level, sheet_type, test_style, created_by, created_at
		FROM practice_sheets
		WHERE course_id = $1
		ORDER BY created_at DESC`

	args := []interface{}{filter.CourseID}
	argIndex := 2

	if filter.Limit > 0 {
		query += fmt.Sprintf(` LIMIT $%d`, argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(` OFFSET $%d`, argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sheets []domain.PracticeSheet
	for rows.Next() {
		var ps domain.PracticeSheet
		if err := rows.Scan(&ps.ID, &ps.CourseID, &ps.TopicID, &ps.StrategyID, &ps.Title, &ps.Level, &ps.SheetType, &ps.TestStyle, &ps.CreatedBy, &ps.CreatedAt); err != nil {
			return nil, err
		}
		sheets = append(sheets, ps)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	exQuery := `
		SELECT pse.id, pse.practice_sheet_id, pse.order_index,
		       e.id, e.topic_id, COALESCE(e.material_id::text,''), e.type, e.question,
		       COALESCE(e.correct_answer,''), COALESCE(e.explanation,''), e.difficulty, e.metadata::text, e.created_at
		FROM practice_sheet_exercises pse
		JOIN exercises e ON e.id = pse.exercise_id
		WHERE pse.practice_sheet_id = ANY($1)
		ORDER BY pse.practice_sheet_id, pse.order_index ASC
	`
	ids := make([]string, len(sheets))
	for i, s := range sheets {
		ids[i] = s.ID
	}
	if len(ids) == 0 {
		return sheets, nil
	}

	exRows, err := r.db.QueryContext(ctx, exQuery, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer exRows.Close()

	index := make(map[string]int, len(sheets))
	for i, s := range sheets {
		index[s.ID] = i
	}
	for exRows.Next() {
		var pse domain.PracticeSheetExercise
		if err := exRows.Scan(
			&pse.ID, &pse.PracticeSheetID, &pse.OrderIndex,
			&pse.Exercise.ID, &pse.Exercise.TopicID, &pse.Exercise.MaterialID, &pse.Exercise.Type,
			&pse.Exercise.Question, &pse.Exercise.CorrectAnswer, &pse.Exercise.Explanation,
			&pse.Exercise.Difficulty, &pse.Exercise.Metadata, &pse.Exercise.CreatedAt,
		); err != nil {
			return nil, err
		}
		if i, ok := index[pse.PracticeSheetID]; ok {
			sheets[i].Exercises = append(sheets[i].Exercises, pse)
		}
	}
	return sheets, nil
}
