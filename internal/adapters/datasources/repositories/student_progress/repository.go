package studentprogress

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Upsert(context.Context, domain.StudentTopicProgress) error
	ListByStudent(context.Context, string) ([]domain.StudentTopicProgress, error)
	ListByStudentAndCourse(ctx context.Context, studentID, courseID string) ([]domain.StudentTopicProgress, error)
	Get(ctx context.Context, studentID, topicID string) (*domain.StudentTopicProgress, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Upsert(ctx context.Context, p domain.StudentTopicProgress) error {
	query := `
		INSERT INTO student_topic_progress (student_id, topic_id, strategy_id, mastery_score, current_level, total_attempts, correct_attempts, streak_days, last_practiced_at, updated_at)
		VALUES ($1, $2, NULLIF($3,'')::uuid, $4, $5, $6, $7, $8, NOW(), NOW())
		ON CONFLICT (student_id, topic_id, strategy_id) DO UPDATE SET
			mastery_score = $4,
			current_level = $5,
			total_attempts = $6,
			correct_attempts = $7,
			streak_days = $8,
			last_practiced_at = NOW(),
			updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, p.StudentID, p.TopicID, p.StrategyID, p.MasteryScore, p.CurrentLevel, p.TotalAttempts, p.CorrectAttempts, p.StreakDays)
	return err
}

func (r *repository) Get(ctx context.Context, studentID, topicID string) (*domain.StudentTopicProgress, error) {
	query := `
		SELECT stp.id, stp.student_id, stp.topic_id, COALESCE(stp.strategy_id::text,''),
		       stp.mastery_score, stp.current_level, stp.total_attempts, stp.correct_attempts,
		       stp.streak_days, stp.last_practiced_at, stp.updated_at, COALESCE(t.title,'')
		FROM student_topic_progress stp
		LEFT JOIN topics t ON t.id = stp.topic_id
		WHERE stp.student_id = $1 AND stp.topic_id = $2::uuid
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, studentID, topicID)
	var p domain.StudentTopicProgress
	err := row.Scan(&p.ID, &p.StudentID, &p.TopicID, &p.StrategyID, &p.MasteryScore, &p.CurrentLevel,
		&p.TotalAttempts, &p.CorrectAttempts, &p.StreakDays, &p.LastPracticedAt, &p.UpdatedAt, &p.TopicTitle)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (r *repository) ListByStudent(ctx context.Context, studentID string) ([]domain.StudentTopicProgress, error) {
	return r.listByFilter(ctx, studentID, "")
}

func (r *repository) ListByStudentAndCourse(ctx context.Context, studentID, courseID string) ([]domain.StudentTopicProgress, error) {
	return r.listByFilter(ctx, studentID, courseID)
}

func (r *repository) listByFilter(ctx context.Context, studentID, courseID string) ([]domain.StudentTopicProgress, error) {
	query := `
		SELECT stp.id, stp.student_id, stp.topic_id, COALESCE(stp.strategy_id::text,''),
		       stp.mastery_score, stp.current_level, stp.total_attempts, stp.correct_attempts,
		       stp.streak_days, stp.last_practiced_at, stp.updated_at, COALESCE(t.title,'')
		FROM student_topic_progress stp
		LEFT JOIN topics t ON t.id = stp.topic_id
		WHERE stp.student_id = $1
	`
	args := []interface{}{studentID}
	if courseID != "" {
		query += ` AND t.course_id = $2::uuid`
		args = append(args, courseID)
	}
	query += ` ORDER BY stp.updated_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.StudentTopicProgress
	for rows.Next() {
		var p domain.StudentTopicProgress
		if err := rows.Scan(&p.ID, &p.StudentID, &p.TopicID, &p.StrategyID, &p.MasteryScore, &p.CurrentLevel,
			&p.TotalAttempts, &p.CorrectAttempts, &p.StreakDays, &p.LastPracticedAt, &p.UpdatedAt, &p.TopicTitle); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}
