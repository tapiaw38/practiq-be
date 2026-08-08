package domain

import "time"

type PracticeSheet struct {
	ID         string
	CourseID   string
	TopicID    string
	StrategyID string
	Title      string
	Level      int
	SheetType  string // 'practice' | 'level_test'
	TestStyle  string // 'keyboard' | 'canvas'
	// ScheduledAt is when students may start the sheet. nil means no schedule.
	ScheduledAt *time.Time
	CreatedBy   string
	CreatedAt   time.Time
	Exercises   []PracticeSheetExercise
}

type PracticeSheetExercise struct {
	ID              string
	PracticeSheetID string
	Exercise        Exercise
	OrderIndex      int
}

type StudentAttempt struct {
	ID              string
	StudentID       string
	ExerciseID      string
	PracticeSheetID string
	AnswerText      string
	ImageURL        string
	AIFeedback      string
	IsCorrect       bool
	Score           float64
	TimeSpentSecs   int
	HintsUsed       int
	// Attachment answers: the uploaded file plus what it is.
	AttachmentURL         string
	AttachmentName        string
	AttachmentContentType string
	// NeedsTeacherReview is set when the assistant could not grade the file.
	NeedsTeacherReview bool
	TeacherIsCorrect   *bool
	TeacherFeedback    string
	TeacherReviewedAt  *time.Time
	CreatedAt          time.Time
}

// PendingAttemptReview is an attachment answer waiting for a teacher, joined
// with the context the teacher needs to judge it.
type PendingAttemptReview struct {
	AttemptID             string
	StudentID             string
	StudentName           string
	ExerciseID            string
	Question              string
	ExerciseType          string
	PracticeSheetID       string
	PracticeSheetTitle    string
	SheetType             string
	CourseID              string
	CourseTitle           string
	AttachmentURL         string
	AttachmentName        string
	AttachmentContentType string
	AnswerText            string
	AIFeedback            string
	TeacherIsCorrect      *bool
	TeacherFeedback       string
	TeacherReviewedAt     *time.Time
	CreatedAt             time.Time
}
