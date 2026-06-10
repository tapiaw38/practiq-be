package domain

import "time"

type LearningStrategy struct {
	ID          string
	Name        string
	Code        string
	Description string
	Status      string
	CreatedAt   time.Time
}

type CourseLearningStrategy struct {
	ID                  string
	CourseID            string
	StrategyID          string
	IsDefault           bool
	Config              string // JSON string
	CreatedAt           time.Time
	StrategyName        string
	StrategyCode        string
	StrategyDescription string
}
