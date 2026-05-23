package subject

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Subject) (string, error)
	List(context.Context) ([]domain.Subject, error)
	Get(context.Context, string) (*domain.Subject, error)
	Update(context.Context, string, domain.Subject) error
	Delete(context.Context, string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, subject domain.Subject) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO subjects (name, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING id
	`, subject.Name, subject.Description, subject.CreatedBy).Scan(&id)
	return id, err
}

func (r *repository) List(ctx context.Context) ([]domain.Subject, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, COALESCE(description, ''), created_by, created_at
		FROM subjects
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subjects := []domain.Subject{}
	for rows.Next() {
		var subject domain.Subject
		if err := rows.Scan(&subject.ID, &subject.Name, &subject.Description, &subject.CreatedBy, &subject.CreatedAt); err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	return subjects, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.Subject, error) {
	var subject domain.Subject
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, COALESCE(description, ''), created_by, created_at
		FROM subjects
		WHERE id = $1
	`, id).Scan(&subject.ID, &subject.Name, &subject.Description, &subject.CreatedBy, &subject.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

func (r *repository) Update(ctx context.Context, id string, subject domain.Subject) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE subjects SET name = $1, description = $2 WHERE id = $3
	`, subject.Name, subject.Description, id)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM subjects WHERE id = $1`, id)
	return err
}
