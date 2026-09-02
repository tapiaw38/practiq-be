package studentprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) ListByStudent(ctx context.Context, studentID string) ([]domain.StudentTopicProgress, error) {
	return r.listByFilter(ctx, studentID, "")
}
