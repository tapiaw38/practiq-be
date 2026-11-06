package courselevel

import "github.com/tapiaw38/practiq-be/internal/domain"

type (
	SheetData struct {
		ID      string `json:"id"`
		Title   string `json:"title"`
		TopicID string `json:"topic_id,omitempty"`

		TopicTitle string `json:"topic_title,omitempty"`
		TopicOrder int    `json:"topic_order"`
		Level      int    `json:"level"`
		SheetType  string `json:"sheet_type"`
		TestStyle  string `json:"test_style"`

		ScheduledAt string `json:"scheduled_at,omitempty"`

		AvailableUntil string `json:"available_until,omitempty"`
		Exercises      int    `json:"exercises"`
		Submitted      bool   `json:"submitted"`

		AttemptsUsed    int   `json:"attempts_used"`
		AttemptsAllowed int   `json:"attempts_allowed"`
		PendingReview   bool  `json:"pending_review"`
		Score           *int  `json:"score,omitempty"`
		Passed          *bool `json:"passed,omitempty"`
	}

	NotebookData struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		TopicID     string `json:"topic_id,omitempty"`
		TopicTitle  string `json:"topic_title,omitempty"`
		TopicOrder  int    `json:"topic_order"`
		Description string `json:"description"`
		Level       int    `json:"level"`
		Pages       int    `json:"pages"`
	}

	LevelData struct {
		Level     int            `json:"level"`
		Unlocked  bool           `json:"unlocked"`
		Practices []SheetData    `json:"practices"`
		LevelTest *SheetData     `json:"level_test"`
		Notebooks []NotebookData `json:"notebooks"`
	}
)

func toSheetData(s domain.PracticeSheet, topicTitle string, topicOrder int) SheetData {
	data := SheetData{
		ID:         s.ID,
		Title:      s.Title,
		TopicID:    s.TopicID,
		TopicTitle: topicTitle,
		TopicOrder: topicOrder,
		Level:      s.Level,
		SheetType:  s.SheetType,
		TestStyle:  s.TestStyle,
		Exercises:  len(s.Exercises),
	}
	if s.ScheduledAt != nil {
		data.ScheduledAt = s.ScheduledAt.UTC().Format("2006-01-02T15:04:05Z")
	}
	if s.AvailableUntil != nil {
		data.AvailableUntil = s.AvailableUntil.UTC().Format("2006-01-02T15:04:05Z")
	}
	return data
}

func toNotebookData(nb domain.Notebook, topicTitle string, topicOrder int) NotebookData {
	return NotebookData{
		ID:          nb.ID,
		Title:       nb.Title,
		TopicID:     nb.TopicID,
		TopicTitle:  topicTitle,
		TopicOrder:  topicOrder,
		Description: nb.Description,
		Level:       nb.Level,
		Pages:       len(nb.Pages),
	}
}
