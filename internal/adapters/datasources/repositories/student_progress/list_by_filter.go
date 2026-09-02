package studentprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) listByFilter(ctx context.Context, studentID, courseID string) ([]domain.StudentTopicProgress, error) {
	query := `
		SELECT stp.id, stp.student_id, stp.topic_id, COALESCE(stp.strategy_id::text,''),
		       stp.mastery_score, stp.current_level, stp.total_attempts, stp.correct_attempts,
		       stp.streak_days, stp.last_practiced_at, stp.updated_at, COALESCE(t.title,'')
		FROM student_topic_progress stp
		LEFT JOIN topics t ON t.id = stp.topic_id
		LEFT JOIN courses c ON c.id = t.course_id
		WHERE stp.student_id = $1 AND (c.id IS NULL OR c.deleted_at IS NULL)
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
