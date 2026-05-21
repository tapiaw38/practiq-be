package domain

import "time"

type StudentTopicProgress struct {
	ID              string
	StudentID       string
	TopicID         string
	StrategyID      string
	MasteryScore    float64
	CurrentLevel    int
	TotalAttempts   int
	CorrectAttempts int
	StreakDays      int
	LastPracticedAt *time.Time
	UpdatedAt       time.Time
	TopicTitle      string // joined field
}
