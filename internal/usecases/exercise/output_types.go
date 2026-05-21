package exercise

import "github.com/tapiaw38/practiq-be/internal/domain"

type ExerciseData struct {
	ID            string `json:"id"`
	TopicID       string `json:"topic_id"`
	MaterialID    string `json:"material_id,omitempty"`
	Type          string `json:"type"`
	Question      string `json:"question"`
	CorrectAnswer string `json:"correct_answer,omitempty"`
	Explanation   string `json:"explanation,omitempty"`
	Difficulty    int    `json:"difficulty"`
	Metadata      string `json:"metadata"`
	CreatedAt     string `json:"created_at"`
}

type ExerciseOutput struct {
	Data ExerciseData `json:"data"`
}

type ExerciseListOutput struct {
	Data []ExerciseData `json:"data"`
}

func toExerciseData(e domain.Exercise) ExerciseData {
	return ExerciseData{
		ID:            e.ID,
		TopicID:       e.TopicID,
		MaterialID:    e.MaterialID,
		Type:          e.Type,
		Question:      e.Question,
		CorrectAnswer: e.CorrectAnswer,
		Explanation:   e.Explanation,
		Difficulty:    e.Difficulty,
		Metadata:      e.Metadata,
		CreatedAt:     e.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
