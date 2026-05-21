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
