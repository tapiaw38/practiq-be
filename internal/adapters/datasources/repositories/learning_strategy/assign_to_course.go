package learningstrategy

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) AssignToCourse(ctx context.Context, cls domain.CourseLearningStrategy) (string, error) {
	query := `
		INSERT INTO course_learning_strategies (course_id, strategy_id, is_default, config, school_id)
		SELECT $1, $2, $3, $4::jsonb, NULLIF($5,'')::uuid
		WHERE EXISTS (SELECT 1 FROM courses WHERE id = $1 AND ($5 = '' OR school_id = NULLIF($5,'')::uuid))
		ON CONFLICT (course_id, strategy_id) DO UPDATE SET is_default = EXCLUDED.is_default, config = EXCLUDED.config
		RETURNING id
	`
	config := cls.Config
	if config == "" {
		config = "{}"
	}
	var id string
	err := r.db.QueryRowContext(ctx, query, cls.CourseID, cls.StrategyID, cls.IsDefault, config, tenant.SchoolID(ctx)).Scan(&id)
	return id, err
}
