package notebook

import (
	"context"
	"fmt"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

type (
	TeacherReviewSubmissionUsecase interface {
		Execute(ctx context.Context, submissionID string, input TeacherReviewInput) (*TeacherReviewSubmissionOutput, error)
	}

	TeacherReviewInput struct {
		IsCorrect bool
		Feedback  string
		TeacherID string
	}

	TeacherReviewSubmissionOutput struct {
		Data NotebookSubmissionFullData `json:"data"`
	}

	teacherReviewSubmissionUsecase struct{ contextFactory appcontext.Factory }
)

func NewTeacherReviewSubmissionUsecase(contextFactory appcontext.Factory) TeacherReviewSubmissionUsecase {
	return &teacherReviewSubmissionUsecase{contextFactory: contextFactory}
}

func (u *teacherReviewSubmissionUsecase) Execute(ctx context.Context, submissionID string, input TeacherReviewInput) (*TeacherReviewSubmissionOutput, error) {
	app := u.contextFactory()
	submission, err := app.Repositories.Notebook.GetFullSubmissionByID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if submission == nil || (input.TeacherID != "" && submission.TeacherID != input.TeacherID) {
		return nil, fmt.Errorf("submission not found")
	}
	if err := app.Repositories.Notebook.UpdateSubmissionTeacherReview(ctx, submissionID, input.IsCorrect, strings.TrimSpace(input.Feedback)); err != nil {
		return nil, err
	}
	updated, err := app.Repositories.Notebook.GetFullSubmissionByID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, fmt.Errorf("submission not found")
	}
	return &TeacherReviewSubmissionOutput{Data: toFullSubmissionData(*updated)}, nil
}
