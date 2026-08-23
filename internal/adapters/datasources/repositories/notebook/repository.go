package notebook

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, n domain.Notebook) (string, error)
	List(ctx context.Context, courseID string) ([]domain.Notebook, error)
	Get(ctx context.Context, id string) (*domain.Notebook, error)
	Update(ctx context.Context, id string, n domain.Notebook) error
	Delete(ctx context.Context, id string) error

	CreatePage(ctx context.Context, p domain.NotebookPage) (string, error)
	GetPage(ctx context.Context, pageID string) (*domain.NotebookPage, error)
	UpdatePage(ctx context.Context, p domain.NotebookPage) error
	DeletePage(ctx context.Context, pageID string) error

	UpsertSubmission(ctx context.Context, s domain.NotebookSubmission) error
	GetSubmission(ctx context.Context, pageID, studentID string) (*domain.NotebookSubmission, error)
	GetSubmissionByID(ctx context.Context, id string) (*domain.NotebookSubmission, error)
	GetFullSubmissionByID(ctx context.Context, id string) (*domain.NotebookSubmissionFull, error)
	ListSubmissions(ctx context.Context, filter SubmissionFilter) ([]domain.NotebookSubmissionFull, error)
	UpdateSubmissionAIReview(ctx context.Context, id string, recognizedText string, isCorrect *bool, feedback string, needsTeacherReview bool) error
	UpdateSubmissionTeacherReview(ctx context.Context, id string, isCorrect bool, feedback string) error
}
type SubmissionFilter struct {
	NotebookID string
	StudentID  string
	CourseID   string
	GradeID    string
	SubjectID  string
	Reviewed   string
	TeacherID  string
	Limit      int
	Offset     int
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
