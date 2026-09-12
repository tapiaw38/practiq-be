package studentprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Upsert(ctx context.Context, p domain.StudentTopicProgress) error {
	query := `
		INSERT INTO student_topic_progress (student_id, topic_id, strategy_id, mastery_score, current_level, total_attempts, correct_attempts, streak_days, last_practiced_at, updated_at)
		VALUES ($1, $2, NULLIF($3,'')::uuid, $4, $5, $6, $7, $8, NOW(), NOW())
		ON CONFLICT (student_id, topic_id) DO UPDATE SET
			strategy_id = COALESCE(NULLIF($3,'')::uuid, student_topic_progress.strategy_id),
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
