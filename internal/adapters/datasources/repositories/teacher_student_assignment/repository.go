package teacherstudentassignment

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Assign(context.Context, domain.TeacherStudentAssignment) error
	Unassign(context.Context, string, string) error
	ListTeachers(context.Context, string) ([]domain.UserProfile, error)
	ListStudents(context.Context, string) ([]domain.UserProfile, error)
	HasAccess(context.Context, string, string) (bool, error)
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
