package notebook

import (
	"context"

	courseRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	GetUsecase interface {
		Execute(ctx context.Context, requesterID string, isAdmin bool, id, studentID string) (*GetOutput, apperrors.ApplicationError)
	}

	GetOutput struct {
		Data NotebookData `json:"data"`
	}

	getUsecase struct{ contextFactory appcontext.Factory }
)

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{contextFactory: contextFactory}
}

func (u *getUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, id, studentID string) (*GetOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	nb, err := app.Repositories.Notebook.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}
	if nb == nil {
		return nil, nil
	}

	if studentID == "" {
		if !isAdmin && nb.TeacherID != requesterID {
			return nil, apperrors.NewForbiddenError()
		}
	} else if !isAdmin {
		if studentID == requesterID {
			hasAccess, err := studentHasNotebookCourseAccess(ctx, app, studentID, nb.CourseID)
			if err != nil {
				return nil, apperrors.NewApplicationError(mappings.NotebookGetError, err)
			}
			if !hasAccess {
				return nil, apperrors.NewForbiddenError()
			}
		} else {
			hasStudentAccess, err := app.Repositories.TeacherStudentAssignment.HasAccess(ctx, requesterID, studentID)
			if err != nil {
				return nil, apperrors.NewApplicationError(mappings.AssignmentListError, err)
			}
			if !hasStudentAccess || nb.TeacherID != requesterID {
				return nil, apperrors.NewForbiddenError()
			}
		}
	}

	// Attach student submissions to each page
	if studentID != "" {
		for i := range nb.Pages {
			sub, _ := app.Repositories.Notebook.GetSubmission(ctx, nb.Pages[i].ID, studentID)
			nb.Pages[i].Submission = sub
		}
	}
	resolveNotebookImages(ctx, app, nb)

	return &GetOutput{Data: toNotebookData(nb)}, nil
}

func studentHasNotebookCourseAccess(ctx context.Context, app *appcontext.Context, studentID, courseID string) (bool, error) {
	courses, err := app.Repositories.Course.List(ctx, courseRepo.ListFilterOptions{StudentID: studentID})
	if err != nil {
		return false, err
	}
	for _, course := range courses {
		if course.ID == courseID {
			return true, nil
		}
	}
	return false, nil
}
