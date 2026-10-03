package studentprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) ListByStudentAndCourse(ctx context.Context, studentID, courseID string) ([]domain.StudentTopicProgress, error) {
	return r.listByFilter(ctx, studentID, courseID)
}
