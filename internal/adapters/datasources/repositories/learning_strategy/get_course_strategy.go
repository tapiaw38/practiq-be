package learningstrategy

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) GetCourseStrategy(ctx context.Context, id string) (*domain.CourseLearningStrategy, error) {
	query := `
		SELECT cls.id, cls.course_id, cls.strategy_id, cls.is_default, COALESCE(cls.config::text, '{}'), cls.created_at,
		       ls.name, ls.code, COALESCE(ls.description,'')
		FROM course_learning_strategies cls
		JOIN learning_strategies ls ON ls.id = cls.strategy_id
		WHERE cls.id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	var cls domain.CourseLearningStrategy
	var configJSON string
	err := row.Scan(&cls.ID, &cls.CourseID, &cls.StrategyID, &cls.IsDefault, &configJSON, &cls.CreatedAt,
		&cls.StrategyName, &cls.StrategyCode, &cls.StrategyDescription)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	cls.Config = configJSON
	return &cls, nil
}
