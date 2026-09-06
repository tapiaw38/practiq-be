package exercise

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.Exercise, error) {
	query := `SELECT e.id, e.topic_id, COALESCE(e.material_id::text,''), e.type, e.question, COALESCE(e.correct_answer,''), COALESCE(e.explanation,''), e.difficulty, e.metadata::text, e.created_at FROM exercises e JOIN topics t ON t.id = e.topic_id JOIN courses c ON c.id = t.course_id WHERE e.id = $1 AND ($2 = '' OR c.school_id = NULLIF($2, '')::uuid)`
	row := r.db.QueryRowContext(ctx, query, id, tenant.SchoolID(ctx))
	var e domain.Exercise
	err := row.Scan(&e.ID, &e.TopicID, &e.MaterialID, &e.Type, &e.Question, &e.CorrectAnswer, &e.Explanation, &e.Difficulty, &e.Metadata, &e.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &e, err
}
