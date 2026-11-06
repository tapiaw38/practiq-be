package domain

import "time"

type PracticeSheet struct {
	ID         string
	CourseID   string
	TopicID    string
	StrategyID string
	Title      string
	Level      int
	SheetType  string
	TestStyle  string

	ScheduledAt *time.Time

	AvailableUntil *time.Time

	MaxAttempts *int

	TimeLimitMinutes *int
	CreatedBy        string
	CreatedAt        time.Time
	Exercises        []PracticeSheetExercise
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

	AttachmentURL         string
	AttachmentName        string
	AttachmentContentType string

	NotGraded bool

	NeedsTeacherReview bool

	AIIsCorrect       *bool
	TeacherIsCorrect  *bool
	TeacherFeedback   string
	TeacherReviewedAt *time.Time
	CreatedAt         time.Time
}

type PendingAttemptReview struct {
	AttemptID         string
	StudentID         string
	StudentName       string
	ExerciseID        string
	Question          string
	ExerciseType      string
	StatementMediaURL string

	HasTeacherImage    bool
	PracticeSheetID    string
	PracticeSheetTitle string
	SheetType          string
	CourseID           string
	CourseTitle        string

	ImageURL              string
	AttachmentURL         string
	AttachmentName        string
	AttachmentContentType string
	AnswerText            string
	AIFeedback            string

	AIIsCorrect       *bool
	TeacherIsCorrect  *bool
	TeacherFeedback   string
	TeacherReviewedAt *time.Time
	CreatedAt         time.Time
}

type SheetOutcome struct {
	Total   int
	Correct int

	Pending int
}

type AttemptContext struct {
	StudentID       string
	PracticeSheetID string
	SheetType       string
	SheetTopicID    string
	ExerciseTopicID string
}

func (p PracticeSheet) AttemptsAllowed() int {
	if p.SheetType != "level_test" {
		return 0
	}
	if p.MaxAttempts != nil && *p.MaxAttempts > 0 {
		return *p.MaxAttempts
	}
	return 1
}

func (p PracticeSheet) Deadline(startedAt *time.Time) *time.Time {
	if p.TimeLimitMinutes == nil || *p.TimeLimitMinutes <= 0 || startedAt == nil {
		return nil
	}
	deadline := startedAt.Add(time.Duration(*p.TimeLimitMinutes) * time.Minute)
	return &deadline
}
