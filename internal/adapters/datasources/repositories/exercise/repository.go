package exercise

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Exercise) (string, error)
	Get(context.Context, string) (*domain.Exercise, error)
	List(context.Context, string) ([]domain.Exercise, error)
	Update(context.Context, string, domain.Exercise) error
	Delete(context.Context, string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, e domain.Exercise) (string, error) {
	query := `
		INSERT INTO exercises (topic_id, type, question, correct_answer, explanation, difficulty, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id string
	metadata := e.Metadata
	if metadata == "" {
		metadata = "{}"
	}
	err := r.db.QueryRowContext(ctx, query, e.TopicID, e.Type, e.Question, e.CorrectAnswer, e.Explanation, e.Difficulty, metadata).Scan(&id)
	return id, err
}

func (r *repository) Get(ctx context.Context, id string) (*domain.Exercise, error) {
	query := `SELECT id, topic_id, COALESCE(material_id::text,''), type, question, COALESCE(correct_answer,''), COALESCE(explanation,''), difficulty, metadata::text, created_at FROM exercises WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var e domain.Exercise
	err := row.Scan(&e.ID, &e.TopicID, &e.MaterialID, &e.Type, &e.Question, &e.CorrectAnswer, &e.Explanation, &e.Difficulty, &e.Metadata, &e.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &e, err
}

func (r *repository) List(ctx context.Context, topicID string) ([]domain.Exercise, error) {
	query := `SELECT id, topic_id, COALESCE(material_id::text,''), type, question, COALESCE(correct_answer,''), COALESCE(explanation,''), difficulty, metadata::text, created_at FROM exercises WHERE topic_id = $1 ORDER BY created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []domain.Exercise
	for rows.Next() {
		var e domain.Exercise
		if err := rows.Scan(&e.ID, &e.TopicID, &e.MaterialID, &e.Type, &e.Question, &e.CorrectAnswer, &e.Explanation, &e.Difficulty, &e.Metadata, &e.CreatedAt); err != nil {
			return nil, err
		}
		exercises = append(exercises, e)
	}
	return exercises, nil
}

func (r *repository) Update(ctx context.Context, id string, e domain.Exercise) error {
	query := `UPDATE exercises SET type=$1, question=$2, correct_answer=$3, explanation=$4, difficulty=$5 WHERE id=$6`
	_, err := r.db.ExecContext(ctx, query, e.Type, e.Question, e.CorrectAnswer, e.Explanation, e.Difficulty, id)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM exercises WHERE id=$1`, id)
	return err
}
