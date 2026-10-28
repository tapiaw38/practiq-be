package learningstrategy

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) AssignToCourse(ctx context.Context, cls domain.CourseLearningStrategy) (string, error) {
	query := `
		INSERT INTO course_learning_strategies (course_id, strategy_id, is_default, config)
		VALUES ($1, $2, $3, $4::jsonb)
		ON CONFLICT (course_id, strategy_id) DO UPDATE SET is_default = EXCLUDED.is_default, config = EXCLUDED.config
		RETURNING id
	`
	config := cls.Config
	if config == "" {
		config = "{}"
	}
	var id string
	err := r.db.QueryRowContext(ctx, query, cls.CourseID, cls.StrategyID, cls.IsDefault, config).Scan(&id)
	return id, err
}
