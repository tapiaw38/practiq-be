package practicesheet

import (
	"context"
	"time"
)

type submissionReceivedAtKey struct{}

func WithSubmissionReceivedAt(ctx context.Context, receivedAt time.Time) context.Context {
	return context.WithValue(ctx, submissionReceivedAtKey{}, receivedAt)
}

func submissionReferenceTime(ctx context.Context) time.Time {
	if receivedAt, ok := ctx.Value(submissionReceivedAtKey{}).(time.Time); ok && !receivedAt.IsZero() {
		return receivedAt
	}
	return time.Now()
}
