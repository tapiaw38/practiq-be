package exercise

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

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
