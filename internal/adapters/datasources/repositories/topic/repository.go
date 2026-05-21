package topic

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Topic) (string, error)
	List(context.Context, string) ([]domain.Topic, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, t domain.Topic) (string, error) {
	query := `
		INSERT INTO topics (course_id, title, description, order_index)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, t.CourseID, t.Title, t.Description, t.OrderIndex).Scan(&id)
	return id, err
}

func (r *repository) List(ctx context.Context, courseID string) ([]domain.Topic, error) {
	query := `SELECT id, course_id, title, description, order_index, created_at FROM topics WHERE course_id = $1 ORDER BY order_index ASC`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []domain.Topic
	for rows.Next() {
		var t domain.Topic
		if err := rows.Scan(&t.ID, &t.CourseID, &t.Title, &t.Description, &t.OrderIndex, &t.CreatedAt); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	return topics, nil
}
