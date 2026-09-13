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
	// AvailableUntil closes the window opened by ScheduledAt. nil means the
	// sheet stays open once it opens.
	AvailableUntil *time.Time
	CreatedBy      string
	CreatedAt      time.Time
	Exercises      []PracticeSheetExercise
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
	// NotGraded means nobody put a verdict on this answer, so it is out of the
	// score rather than counted as wrong. Every reader has to honour it, or the
	// stored is_correct=false reads as a mistake the student made.
	NotGraded bool
	// NeedsTeacherReview says a person still has to settle it. Only a level
	// test sets it; a practice is never queued for anyone.
	NeedsTeacherReview bool
	// AIIsCorrect is the assistant's suggestion, nil when it had none.
	AIIsCorrect       *bool
	TeacherIsCorrect  *bool
	TeacherFeedback   string
	TeacherReviewedAt *time.Time
	CreatedAt         time.Time
}

// PendingAttemptReview is an answer waiting for a teacher, joined with the
// context and statement material the teacher needs to judge it.
type PendingAttemptReview struct {
	AttemptID         string
	StudentID         string
	StudentName       string
	ExerciseID        string
	Question          string
	ExerciseType      string
	StatementMediaURL string
	// HasTeacherImage says a handwritten statement exists without carrying it.
	HasTeacherImage    bool
	PracticeSheetID    string
	PracticeSheetTitle string
	SheetType          string
	CourseID           string
	CourseTitle        string
	// ImageURL is the canvas the student drew, when the answer was handwritten.
	// Without it the teacher grades work they cannot see.
	ImageURL              string
	AttachmentURL         string
	AttachmentName        string
	AttachmentContentType string
	AnswerText            string
	AIFeedback            string
	// AIIsCorrect is what the assistant suggested, for the teacher to confirm.
	AIIsCorrect       *bool
	TeacherIsCorrect  *bool
	TeacherFeedback   string
	TeacherReviewedAt *time.Time
	CreatedAt         time.Time
}

// SheetOutcome is a student's standing on a sheet, counting the latest answer
// per exercise.
type SheetOutcome struct {
	Total   int
	Correct int
	// Pending counts answers still waiting for a teacher.
	Pending int
}

// AttemptContext is what a review needs to know to recompute progress.
type AttemptContext struct {
	StudentID       string
	PracticeSheetID string
	SheetType       string
	SheetTopicID    string
	ExerciseTopicID string
}
