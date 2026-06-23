package practicesheet

import "github.com/tapiaw38/practiq-be/internal/domain"

type (
	ExerciseData struct {
		ID            string `json:"id"`
		TopicID       string `json:"topic_id"`
		Type          string `json:"type"`
		Question      string `json:"question"`
		CorrectAnswer string `json:"correct_answer,omitempty"`
		Explanation   string `json:"explanation,omitempty"`
		Difficulty    int    `json:"difficulty"`
		Metadata      string `json:"metadata"`
	}

	SheetExerciseData struct {
		ID         string       `json:"id"`
		OrderIndex int          `json:"order_index"`
		Exercise   ExerciseData `json:"exercise"`
	}

	PracticeSheetData struct {
		ID         string              `json:"id"`
		CourseID   string              `json:"course_id"`
		TopicID    string              `json:"topic_id"`
		StrategyID string              `json:"strategy_id"`
		Title      string              `json:"title"`
		Level      int                 `json:"level"`
		SheetType  string              `json:"sheet_type"`
		TestStyle  string              `json:"test_style"`
		CreatedBy  string              `json:"created_by"`
		CreatedAt  string              `json:"created_at"`
		Exercises  []SheetExerciseData `json:"exercises"`
	}

	ExerciseResultData struct {
		ExerciseID    string `json:"exercise_id"`
		IsCorrect     bool   `json:"is_correct"`
		StudentAnswer string `json:"student_answer"`
		CorrectAnswer string `json:"correct_answer"`
		AIFeedback    string `json:"ai_feedback,omitempty"`
	}

	SubmitResult struct {
		Score           float64              `json:"score"`
		Correct         int                  `json:"correct"`
		Total           int                  `json:"total"`
		MasteryScore    float64              `json:"mastery_score"`
		Recommendation  string               `json:"recommendation"`
		AIFeedback      string               `json:"ai_feedback,omitempty"`
		ShouldLevelUp   bool                 `json:"should_level_up"`
		ShouldRepeat    bool                 `json:"should_repeat"`
		NextLevel       int                  `json:"next_level"`
		ExerciseResults []ExerciseResultData `json:"exercise_results"`
	}
)

func toSheetData(ps domain.PracticeSheet) PracticeSheetData {
	exercises := make([]SheetExerciseData, 0, len(ps.Exercises))
	for _, pse := range ps.Exercises {
		exercises = append(exercises, SheetExerciseData{
			ID:         pse.ID,
			OrderIndex: pse.OrderIndex,
			Exercise: ExerciseData{
				ID:            pse.Exercise.ID,
				TopicID:       pse.Exercise.TopicID,
				Type:          pse.Exercise.Type,
				Question:      pse.Exercise.Question,
				CorrectAnswer: pse.Exercise.CorrectAnswer,
				Explanation:   pse.Exercise.Explanation,
				Difficulty:    pse.Exercise.Difficulty,
				Metadata:      pse.Exercise.Metadata,
			},
		})
	}
	sheetType := ps.SheetType
	if sheetType == "" {
		sheetType = "practice"
	}
	testStyle := ps.TestStyle
	if testStyle == "" {
		testStyle = "keyboard"
	}
	return PracticeSheetData{
		ID:         ps.ID,
		CourseID:   ps.CourseID,
		TopicID:    ps.TopicID,
		StrategyID: ps.StrategyID,
		Title:      ps.Title,
		Level:      ps.Level,
		SheetType:  sheetType,
		TestStyle:  testStyle,
		CreatedBy:  ps.CreatedBy,
		CreatedAt:  ps.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Exercises:  exercises,
	}
}

func toSubmitOutputData(score float64, correct, total int, masteryScore float64, recommendation, aiFeedback string, shouldLevelUp, shouldRepeat bool, nextLevel int, exerciseResults []ExerciseResultData) SubmitResult {
	return SubmitResult{
		Score:           score,
		Correct:         correct,
		Total:           total,
		MasteryScore:    masteryScore,
		Recommendation:  recommendation,
		AIFeedback:      aiFeedback,
		ShouldLevelUp:   shouldLevelUp,
		ShouldRepeat:    shouldRepeat,
		NextLevel:       nextLevel,
		ExerciseResults: exerciseResults,
	}
}
