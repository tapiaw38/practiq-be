package courselevel

import "github.com/tapiaw38/practiq-be/internal/domain"

type (
	SheetData struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		Level     int    `json:"level"`
		SheetType string `json:"sheet_type"`
		TestStyle string `json:"test_style"`
		// ScheduledAt is UTC and empty when the sheet has no date. The client
		// uses it to show the date and disable the sheet until then.
		ScheduledAt string `json:"scheduled_at,omitempty"`
		// AvailableUntil is UTC and empty when the sheet remains open.
		AvailableUntil string `json:"available_until,omitempty"`
		Exercises      int    `json:"exercises"`
	}

	NotebookData struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
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

func toSheetData(s domain.PracticeSheet) SheetData {
	data := SheetData{
		ID:        s.ID,
		Title:     s.Title,
		Level:     s.Level,
		SheetType: s.SheetType,
		TestStyle: s.TestStyle,
		Exercises: len(s.Exercises),
	}
	if s.ScheduledAt != nil {
		data.ScheduledAt = s.ScheduledAt.UTC().Format("2006-01-02T15:04:05Z")
	}
	if s.AvailableUntil != nil {
		data.AvailableUntil = s.AvailableUntil.UTC().Format("2006-01-02T15:04:05Z")
	}
	return data
}

func toNotebookData(nb domain.Notebook) NotebookData {
	return NotebookData{
		ID:          nb.ID,
		Title:       nb.Title,
		Description: nb.Description,
		Level:       nb.Level,
		Pages:       len(nb.Pages),
	}
}
