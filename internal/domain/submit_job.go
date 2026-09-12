package domain

import "time"

type SubmitJob struct {
	ID        string
	Kind      string // "practice_sheet" or "notebook"
	StudentID string
	Status    string
	ErrorCode string
	Message   string
	Result    []byte // JSON
	CreatedAt time.Time
	UpdatedAt time.Time
}
