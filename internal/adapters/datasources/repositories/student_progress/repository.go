package studentprogress

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Upsert(context.Context, domain.StudentTopicProgress) error
	ListByStudent(context.Context, string) ([]domain.StudentTopicProgress, error)
	ListByStudentAndCourse(ctx context.Context, studentID, courseID string) ([]domain.StudentTopicProgress, error)
	Get(ctx context.Context, studentID, topicID string) (*domain.StudentTopicProgress, error)
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
