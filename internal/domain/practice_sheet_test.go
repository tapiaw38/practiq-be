package domain

import (
	"testing"
	"time"
)

func intPtr(v int) *int { return &v }

func TestAttemptsAllowed(t *testing.T) {

	if got := (PracticeSheet{SheetType: "practice", MaxAttempts: intPtr(3)}).AttemptsAllowed(); got != 0 {
		t.Fatalf("a practice allows %d attempts, want it uncounted", got)
	}

	if got := (PracticeSheet{SheetType: "level_test"}).AttemptsAllowed(); got != 1 {
		t.Fatalf("a level test with no limit allows %d, want 1", got)
	}
	if got := (PracticeSheet{SheetType: "level_test", MaxAttempts: intPtr(3)}).AttemptsAllowed(); got != 3 {
		t.Fatalf("AttemptsAllowed() = %d, want the limit the teacher set", got)
	}

	if got := (PracticeSheet{SheetType: "level_test", MaxAttempts: intPtr(0)}).AttemptsAllowed(); got != 1 {
		t.Fatalf("AttemptsAllowed() = %d for a zero limit, want the default of 1", got)
	}
}

func TestDeadline(t *testing.T) {
	started := time.Date(2026, 9, 16, 10, 0, 0, 0, time.UTC)

	sheet := PracticeSheet{SheetType: "level_test", TimeLimitMinutes: intPtr(45)}
	deadline := sheet.Deadline(&started)
	if deadline == nil || !deadline.Equal(started.Add(45*time.Minute)) {
		t.Fatalf("Deadline() = %v, want 45 minutes after the start", deadline)
	}

	if got := (PracticeSheet{SheetType: "level_test"}).Deadline(&started); got != nil {
		t.Fatalf("Deadline() = %v with no limit, want none", got)
	}
	if got := sheet.Deadline(nil); got != nil {
		t.Fatalf("Deadline() = %v for a student who never started, want none", got)
	}
}
