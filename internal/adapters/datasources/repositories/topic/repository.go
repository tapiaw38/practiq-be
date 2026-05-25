package topic

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Topic) (string, error)
	List(context.Context, string) ([]domain.Topic, error)
	Get(context.Context, string) (*domain.Topic, error)
	Update(context.Context, string, domain.Topic) error
	Delete(context.Context, string) error
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

func (r *repository) Get(ctx context.Context, id string) (*domain.Topic, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, course_id, title, description, order_index, created_at
		FROM topics WHERE id = $1
	`, id)

	var t domain.Topic
	if err := row.Scan(&t.ID, &t.CourseID, &t.Title, &t.Description, &t.OrderIndex, &t.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *repository) Update(ctx context.Context, id string, t domain.Topic) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE topics SET title = $1, description = $2, order_index = $3 WHERE id = $4
	`, t.Title, t.Description, t.OrderIndex, id)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM topics WHERE id = $1`, id)
	return err
}
