package studentprogress

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Get(ctx context.Context, studentID, topicID string) (*domain.StudentTopicProgress, error) {
	query := `
		SELECT stp.id, stp.student_id, stp.topic_id, COALESCE(stp.strategy_id::text,''),
		       stp.mastery_score, stp.current_level, stp.total_attempts, stp.correct_attempts,
		       stp.streak_days, stp.last_practiced_at, stp.updated_at, COALESCE(t.title,'')
		FROM student_topic_progress stp
		LEFT JOIN topics t ON t.id = stp.topic_id
		WHERE stp.student_id = $1 AND stp.topic_id = $2::uuid AND ($3 = '' OR stp.school_id = NULLIF($3, '')::uuid)
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, studentID, topicID, tenant.SchoolID(ctx))
	var p domain.StudentTopicProgress
	err := row.Scan(&p.ID, &p.StudentID, &p.TopicID, &p.StrategyID, &p.MasteryScore, &p.CurrentLevel,
		&p.TotalAttempts, &p.CorrectAttempts, &p.StreakDays, &p.LastPracticedAt, &p.UpdatedAt, &p.TopicTitle)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}
