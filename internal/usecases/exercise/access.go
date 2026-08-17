package exercise

import (
	"context"

	reposCourse "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

// requesterCanReadExercise allows the course's teacher and anyone taking the
// course, which is the same audience the exercise listing already serves.
func requesterCanReadExercise(ctx context.Context, app *appcontext.Context, requesterID string, isAdmin bool, topicID string) apperrors.ApplicationError {
	if isAdmin {
		return nil
	}
	topic, err := app.Repositories.Topic.Get(ctx, topicID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.TopicGetError, err)
	}
	if topic == nil {
		return apperrors.NewNotFoundError("topic not found")
	}
	course, err := app.Repositories.Course.Get(ctx, topic.CourseID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if course == nil {
		return apperrors.NewNotFoundError("course not found")
	}
	if course.TeacherID == requesterID {
		return nil
	}
	courses, err := app.Repositories.Course.List(ctx, reposCourse.ListFilterOptions{StudentID: requesterID})
	if err != nil {
		return apperrors.NewApplicationError(mappings.CourseListError, err)
	}
	for _, c := range courses {
		if c.ID == topic.CourseID {
			return nil
		}
	}
	return apperrors.NewForbiddenError()
}

func requesterCanWriteTopic(ctx context.Context, app *appcontext.Context, requesterID string, isAdmin bool, topicID string) apperrors.ApplicationError {
	topic, err := app.Repositories.Topic.Get(ctx, topicID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.TopicGetError, err)
	}
	if topic == nil {
		return apperrors.NewNotFoundError("topic not found")
	}
	return requesterCanWriteCourse(ctx, app, requesterID, isAdmin, topic.CourseID)
}

func requesterCanWriteCourse(ctx context.Context, app *appcontext.Context, requesterID string, isAdmin bool, courseID string) apperrors.ApplicationError {
	if isAdmin {
		return nil
	}
	course, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if course == nil {
		return apperrors.NewNotFoundError("course not found")
	}
	if course.TeacherID != requesterID {
		return apperrors.NewForbiddenError()
	}
	return nil
}
