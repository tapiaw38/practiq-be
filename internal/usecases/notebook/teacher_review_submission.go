package notebook

import (
	"context"
	"fmt"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

type (
	TeacherReviewSubmissionUsecase interface {
		Execute(ctx context.Context, submissionID, teacherID, bearerToken string, input TeacherReviewInput) (*TeacherReviewSubmissionOutput, error)
	}

	TeacherReviewInput struct {
		IsCorrect bool   `json:"teacher_is_correct"`
		Feedback  string `json:"teacher_feedback"`
	}

	TeacherReviewSubmissionOutput struct {
		Data NotebookSubmissionFullData `json:"data"`
	}

	teacherReviewSubmissionUsecase struct{ contextFactory appcontext.Factory }
)

func NewTeacherReviewSubmissionUsecase(contextFactory appcontext.Factory) TeacherReviewSubmissionUsecase {
	return &teacherReviewSubmissionUsecase{contextFactory: contextFactory}
}

func (u *teacherReviewSubmissionUsecase) Execute(ctx context.Context, submissionID, teacherID, bearerToken string, input TeacherReviewInput) (*TeacherReviewSubmissionOutput, error) {
	app := u.contextFactory()
	submission, err := app.Repositories.Notebook.GetFullSubmissionByID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if submission == nil || (teacherID != "" && submission.TeacherID != teacherID) {
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

	names, err := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, []string{updated.StudentID})
	if err != nil {
		return nil, err
	}
	info := names[updated.StudentID]
	updated.StudentName = identity.FullName(info, updated.StudentID)
	updated.StudentEmail = info.Email

	return &TeacherReviewSubmissionOutput{Data: toFullSubmissionData(*updated)}, nil
}
