package material

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Material) (string, error)
	List(context.Context, string) ([]domain.Material, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, m domain.Material) (string, error) {
	query := `
		INSERT INTO materials (course_id, teacher_id, title, type, file_url, extracted_text, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, m.CourseID, m.TeacherID, m.Title, m.Type, m.FileURL, m.ExtractedText, m.Status).Scan(&id)
	return id, err
}

func (r *repository) List(ctx context.Context, courseID string) ([]domain.Material, error) {
	query := `SELECT id, course_id, teacher_id, title, type, COALESCE(file_url,''), COALESCE(extracted_text,''), status, created_at FROM materials WHERE course_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []domain.Material
	for rows.Next() {
		var m domain.Material
		if err := rows.Scan(&m.ID, &m.CourseID, &m.TeacherID, &m.Title, &m.Type, &m.FileURL, &m.ExtractedText, &m.Status, &m.CreatedAt); err != nil {
			return nil, err
		}
		materials = append(materials, m)
	}
	return materials, nil
}
