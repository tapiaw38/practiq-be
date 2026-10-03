package domain

import "time"

type SubmitJob struct {
	ID        string
	Kind      string
	StudentID string
	Status    string
	ErrorCode string
	Message   string
	Result    []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}
