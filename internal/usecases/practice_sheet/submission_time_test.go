package practicesheet

import (
	"context"
	"testing"
	"time"
)

func TestSubmissionReferenceTimeUsesAsyncReceiptTime(t *testing.T) {
	receivedAt := time.Date(2026, time.September, 17, 12, 0, 0, 0, time.UTC)
	got := submissionReferenceTime(WithSubmissionReceivedAt(context.Background(), receivedAt))
	if !got.Equal(receivedAt) {
		t.Fatalf("submissionReferenceTime() = %v, want receipt time %v", got, receivedAt)
	}
}

func TestSubmissionReferenceTimeUsesCurrentTimeWithoutAsyncReceipt(t *testing.T) {
	before := time.Now()
	got := submissionReferenceTime(context.Background())
	after := time.Now()
	if got.Before(before) || got.After(after) {
		t.Fatalf("submissionReferenceTime() = %v, want current time between %v and %v", got, before, after)
	}
}
