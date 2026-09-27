package notebook

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

func ensureTopicBelongsToCourse(ctx context.Context, app *appcontext.Context, topicID, courseID string) apperrors.ApplicationError {
	topicID = strings.TrimSpace(topicID)
	if topicID == "" {
		return nil
	}
	topic, err := app.Repositories.Topic.Get(ctx, topicID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.TopicListError, err)
	}
	if topic == nil || topic.CourseID != courseID {
		return apperrors.NewBadRequestError("topic must belong to the notebook course")
	}
	return nil
}
