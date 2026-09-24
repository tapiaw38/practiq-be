package studentattempt

import "context"

func (r *repository) SaveCanvasWork(ctx context.Context, attemptID, imageData string) error {
	// Canvas image location is persisted in student_attempts.image_url.
	// student_work_canvas was removed because it duplicated attempt data.
	return nil
}
