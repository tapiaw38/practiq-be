package notebook

import (
	"context"

	notebookrepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/notebook"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

type (
	ListSubmissionsUsecase interface {
		Execute(ctx context.Context, input ListSubmissionsInput) (*ListSubmissionsOutput, error)
	}

	ListSubmissionsInput struct {
		NotebookID string
		StudentID  string
		CourseID   string
		Reviewed   string
		TeacherID  string
		Limit      int
		Offset     int
	}

	ListSubmissionsOutput struct {
		Data []NotebookSubmissionFullData `json:"data"`
	}

	listSubmissionsUsecase struct{ contextFactory appcontext.Factory }
)

func NewListSubmissionsUsecase(contextFactory appcontext.Factory) ListSubmissionsUsecase {
	return &listSubmissionsUsecase{contextFactory: contextFactory}
}

func (u *listSubmissionsUsecase) Execute(ctx context.Context, input ListSubmissionsInput) (*ListSubmissionsOutput, error) {
	app := u.contextFactory()
	items, err := app.Repositories.Notebook.ListSubmissions(ctx, notebookrepo.SubmissionFilter{
		NotebookID: input.NotebookID,
		StudentID:  input.StudentID,
		CourseID:   input.CourseID,
		Reviewed:   input.Reviewed,
		TeacherID:  input.TeacherID,
		Limit:      input.Limit,
		Offset:     input.Offset,
	})
	if err != nil {
		return nil, err
	}
	data := make([]NotebookSubmissionFullData, 0, len(items))
	for _, item := range items {
		data = append(data, toFullSubmissionData(item))
	}
	return &ListSubmissionsOutput{Data: data}, nil
}
