package domain

import "time"

type Exercise struct {
	ID            string
	TopicID       string
	MaterialID    string
	Type          string
	Question      string
	CorrectAnswer string
	Explanation   string
	Difficulty    int
	Metadata      string // JSON string
	CreatedAt     time.Time
}
