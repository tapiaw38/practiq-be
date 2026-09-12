package practicesheet

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.PracticeSheet) (string, error)
	AddExercise(ctx context.Context, sheetID, exerciseID string, orderIndex int) error
	ReplaceExercises(ctx context.Context, sheetID string, exerciseIDs []string) error
	Get(context.Context, string) (*domain.PracticeSheet, error)
	List(context.Context, ListFilter) ([]domain.PracticeSheet, error)
	Update(context.Context, string, domain.PracticeSheet) error
	Delete(context.Context, string) error
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
