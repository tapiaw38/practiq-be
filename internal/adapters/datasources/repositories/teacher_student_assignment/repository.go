package teacherstudentassignment

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Assign(context.Context, domain.TeacherStudentAssignment) error
	Unassign(context.Context, string, string) error
	ListTeachers(context.Context, ListFilter) ([]domain.UserProfile, error)
	ListStudents(context.Context, ListFilter) ([]domain.UserProfile, error)
	HasAccess(context.Context, string, string) (bool, error)
}
type ListFilter struct {
	UserID string
	Limit  int
	Offset int
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
