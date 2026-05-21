package courselevel

import "github.com/tapiaw38/practiq-be/internal/domain"

type SheetData struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Level     int    `json:"level"`
	SheetType string `json:"sheet_type"`
	TestStyle string `json:"test_style"`
	Exercises int    `json:"exercises"`
}

type NotebookData struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       int    `json:"level"`
	Pages       int    `json:"pages"`
}

type LevelData struct {
	Level     int          `json:"level"`
	Unlocked  bool         `json:"unlocked"`
	Practices []SheetData  `json:"practices"`
	LevelTest *SheetData   `json:"level_test"`
	Notebooks []NotebookData `json:"notebooks"`
}

type CourseLevelsOutput struct {
	CurrentLevel int         `json:"current_level"`
	Levels       []LevelData `json:"levels"`
}

func toSheetData(s domain.PracticeSheet) SheetData {
	return SheetData{
		ID:        s.ID,
		Title:     s.Title,
		Level:     s.Level,
		SheetType: s.SheetType,
		TestStyle: s.TestStyle,
		Exercises: len(s.Exercises),
	}
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
