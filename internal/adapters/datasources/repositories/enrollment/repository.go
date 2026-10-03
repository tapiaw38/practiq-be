package enrollment

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Enrollment) error
	ListStudents(context.Context, ListFilter) ([]domain.UserProfile, error)
	Exists(context.Context, string, string) (bool, error)
}

type ListFilter struct {
	CourseID string
	Limit    int
	Offset   int
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
