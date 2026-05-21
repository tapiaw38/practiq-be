package domain

import "time"

type AIConversation struct {
	ID              string
	StudentID       string
	CourseID        string
	PracticeSheetID string
	CreatedAt       time.Time
}

type AIMessage struct {
	ID             string
	ConversationID string
	Sender         string
	MessageType    string
	Content        string
	CreatedAt      time.Time
}

type AIHelpRequest struct {
	ID         string
	StudentID  string
	ExerciseID string
	Question   string
	AIResponse string
	HelpType   string
	CreatedAt  time.Time
}
