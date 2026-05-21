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
	CreatedBy  string
	CreatedAt  time.Time
	Exercises  []PracticeSheetExercise
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
	IsCorrect       bool
	Score           float64
	TimeSpentSecs   int
	HintsUsed       int
	CreatedAt       time.Time
}
