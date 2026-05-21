package ai

import "github.com/tapiaw38/practiq-be/internal/domain"

type ConversationData struct {
	ID              string `json:"id"`
	StudentID       string `json:"student_id"`
	CourseID        string `json:"course_id"`
	PracticeSheetID string `json:"practice_sheet_id"`
	CreatedAt       string `json:"created_at"`
}

type ConversationOutput struct {
	Data ConversationData `json:"data"`
}

type MessageData struct {
	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`
	Sender         string `json:"sender"`
	MessageType    string `json:"message_type"`
	Content        string `json:"content"`
	CreatedAt      string `json:"created_at"`
}

type MessagesListOutput struct {
	Data []MessageData `json:"data"`
}

type HelpOutput struct {
	Data HelpData `json:"data"`
}

type HelpData struct {
	ID       string `json:"id"`
	Response string `json:"response"`
	HelpType string `json:"help_type"`
}

func toConversationData(c domain.AIConversation) ConversationData {
	return ConversationData{
		ID:              c.ID,
		StudentID:       c.StudentID,
		CourseID:        c.CourseID,
		PracticeSheetID: c.PracticeSheetID,
		CreatedAt:       c.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toMessageData(m domain.AIMessage) MessageData {
	return MessageData{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		Sender:         m.Sender,
		MessageType:    m.MessageType,
		Content:        m.Content,
		CreatedAt:      m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
