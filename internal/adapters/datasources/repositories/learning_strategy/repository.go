package learningstrategy

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	List(context.Context) ([]domain.LearningStrategy, error)
	Get(context.Context, string) (*domain.LearningStrategy, error)
	GetByCode(context.Context, string) (*domain.LearningStrategy, error)
	Create(context.Context, domain.LearningStrategy) (string, error)
	Update(context.Context, string, domain.LearningStrategy) error
	Delete(context.Context, string) error

	// Course learning strategies
	ListByCourse(ctx context.Context, courseID string) ([]domain.CourseLearningStrategy, error)
	AssignToCourse(ctx context.Context, cls domain.CourseLearningStrategy) (string, error)
	GetCourseStrategy(ctx context.Context, id string) (*domain.CourseLearningStrategy, error)
	UnassignFromCourse(ctx context.Context, id string) error
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
