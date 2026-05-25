package practicesheet

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.PracticeSheet) (string, error)
	AddExercise(ctx context.Context, sheetID, exerciseID string, orderIndex int) error
	Get(context.Context, string) (*domain.PracticeSheet, error)
	List(context.Context, string) ([]domain.PracticeSheet, error)
	Update(context.Context, string, domain.PracticeSheet) error
	Delete(context.Context, string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, ps domain.PracticeSheet) (string, error) {
	sheetType := ps.SheetType
	if sheetType != "level_test" {
		sheetType = "practice"
	}
	testStyle := ps.TestStyle
	if testStyle != "canvas" {
		testStyle = "keyboard"
	}
	query := `
		INSERT INTO practice_sheets (course_id, topic_id, strategy_id, title, level, sheet_type, test_style, created_by)
		VALUES ($1, NULLIF($2,'')::uuid, NULLIF($3,'')::uuid, $4, $5, $6, $7, $8)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, ps.CourseID, ps.TopicID, ps.StrategyID, ps.Title, ps.Level, sheetType, testStyle, ps.CreatedBy).Scan(&id)
	return id, err
}

func (r *repository) AddExercise(ctx context.Context, sheetID, exerciseID string, orderIndex int) error {
	query := `INSERT INTO practice_sheet_exercises (practice_sheet_id, exercise_id, order_index) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, sheetID, exerciseID, orderIndex)
	return err
}

func (r *repository) Get(ctx context.Context, id string) (*domain.PracticeSheet, error) {
	query := `
		SELECT ps.id, ps.course_id, COALESCE(ps.topic_id::text,''), COALESCE(ps.strategy_id::text,''), ps.title, ps.level, ps.sheet_type, ps.test_style, ps.created_by, ps.created_at
		FROM practice_sheets ps
		WHERE ps.id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	var ps domain.PracticeSheet
	err := row.Scan(&ps.ID, &ps.CourseID, &ps.TopicID, &ps.StrategyID, &ps.Title, &ps.Level, &ps.SheetType, &ps.TestStyle, &ps.CreatedBy, &ps.CreatedAt)
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
		WHERE pse.practice_sheet_id = $1
		ORDER BY pse.order_index ASC
	`
	rows, err := r.db.QueryContext(ctx, exQuery, id)
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

func (r *repository) Update(ctx context.Context, id string, ps domain.PracticeSheet) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE practice_sheets SET title = $1, topic_id = NULLIF($2,'')::uuid WHERE id = $3
	`, ps.Title, ps.TopicID, id)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM practice_sheets WHERE id = $1`, id)
	return err
}

func (r *repository) List(ctx context.Context, courseID string) ([]domain.PracticeSheet, error) {
	query := `
		SELECT id, course_id, COALESCE(topic_id::text,''), COALESCE(strategy_id::text,''), title, level, sheet_type, test_style, created_by, created_at
		FROM practice_sheets
		WHERE course_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, courseID)
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
	return sheets, nil
}
