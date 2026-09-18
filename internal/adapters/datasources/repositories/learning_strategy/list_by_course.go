package learningstrategy

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) ListByCourse(ctx context.Context, courseID string) ([]domain.CourseLearningStrategy, error) {
	query := `
		SELECT cls.id, cls.course_id, cls.strategy_id, cls.is_default, COALESCE(cls.config::text, '{}'), cls.created_at,
		       ls.name, ls.code, COALESCE(ls.description,'')
		FROM course_learning_strategies cls
		JOIN learning_strategies ls ON ls.id = cls.strategy_id
		WHERE cls.course_id = $1
		ORDER BY cls.created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.CourseLearningStrategy
	for rows.Next() {
		var cls domain.CourseLearningStrategy
		var configJSON string
		if err := rows.Scan(&cls.ID, &cls.CourseID, &cls.StrategyID, &cls.IsDefault, &configJSON, &cls.CreatedAt,
			&cls.StrategyName, &cls.StrategyCode, &cls.StrategyDescription); err != nil {
			return nil, err
		}
		cls.Config = configJSON
		list = append(list, cls)
	}
	return list, nil
}
