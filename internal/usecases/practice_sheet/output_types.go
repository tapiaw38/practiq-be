package practicesheet

import (
	"encoding/json"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

const timeFormat = "2006-01-02T15:04:05Z"

// mediaLinkTTL has to outlive solving the sheet without leaving a shareable
// link around for long. Reloading the sheet issues fresh ones.
const mediaLinkTTL = time.Hour

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
		// MediaViewURL is the temporary URL for the statement's attached media.
		// Metadata keeps the canonical one so it can be written back unchanged.
		MediaViewURL string `json:"media_view_url,omitempty"`
	}

	SheetExerciseData struct {
		ID         string       `json:"id"`
		OrderIndex int          `json:"order_index"`
		Exercise   ExerciseData `json:"exercise"`
	}

	PracticeSheetData struct {
		ID         string `json:"id"`
		CourseID   string `json:"course_id"`
		TopicID    string `json:"topic_id"`
		StrategyID string `json:"strategy_id"`
		Title      string `json:"title"`
		Level      int    `json:"level"`
		SheetType  string `json:"sheet_type"`
		TestStyle  string `json:"test_style"`
		// ScheduledAt is empty when the sheet can be taken at any time.
		ScheduledAt string              `json:"scheduled_at,omitempty"`
		CreatedBy   string              `json:"created_by"`
		CreatedAt   string              `json:"created_at"`
		Exercises   []SheetExerciseData `json:"exercises"`
	}

	ExerciseResultData struct {
		ExerciseID    string `json:"exercise_id"`
		IsCorrect     bool   `json:"is_correct"`
		StudentAnswer string `json:"student_answer"`
		CorrectAnswer string `json:"correct_answer"`
		AIFeedback    string `json:"ai_feedback,omitempty"`
		// NeedsTeacherReview marks an attachment the assistant could not grade.
		NeedsTeacherReview bool `json:"needs_teacher_review,omitempty"`
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
		PendingReview   bool                 `json:"pending_review,omitempty"`
		NextLevel       int                  `json:"next_level"`
		ExerciseResults []ExerciseResultData `json:"exercise_results"`
	}
)

func toSheetData(app *appcontext.Context, ps domain.PracticeSheet, includeTeacherData bool) PracticeSheetData {
	exercises := make([]SheetExerciseData, 0, len(ps.Exercises))
	for _, pse := range ps.Exercises {
		metadata := metadataWithoutMediaURL(pse.Exercise.Metadata)
		if includeTeacherData {
			metadata = pse.Exercise.Metadata
		}
		exercise := ExerciseData{
			ID:           pse.Exercise.ID,
			TopicID:      pse.Exercise.TopicID,
			Type:         pse.Exercise.Type,
			Question:     pse.Exercise.Question,
			Difficulty:   pse.Exercise.Difficulty,
			Metadata:     metadata,
			MediaViewURL: mediaViewURL(app, pse.Exercise),
		}
		// Correct answers and explanations are teacher-only data. The backend
		// grades submissions, so students never need either in the sheet payload.
		if includeTeacherData {
			exercise.CorrectAnswer = pse.Exercise.CorrectAnswer
			exercise.Explanation = pse.Exercise.Explanation
		}
		exercises = append(exercises, SheetExerciseData{
			ID:         pse.ID,
			OrderIndex: pse.OrderIndex,
			Exercise:   exercise,
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
	data := PracticeSheetData{
		ID:         ps.ID,
		CourseID:   ps.CourseID,
		TopicID:    ps.TopicID,
		StrategyID: ps.StrategyID,
		Title:      ps.Title,
		Level:      ps.Level,
		SheetType:  sheetType,
		TestStyle:  testStyle,
		CreatedBy:  ps.CreatedBy,
		CreatedAt:  ps.CreatedAt.Format(timeFormat),
		Exercises:  exercises,
	}
	if ps.ScheduledAt != nil {
		data.ScheduledAt = ps.ScheduledAt.UTC().Format(timeFormat)
	}
	return data
}

// Practice sheets are read by students. Keep non-media exercise metadata
// (fill-blanks options, attachment rules, etc.) but never expose the stable
// bucket URL; MediaViewURL is the short-lived browser URL instead.
func metadataWithoutMediaURL(metadata string) string {
	var values map[string]json.RawMessage
	if err := json.Unmarshal([]byte(metadata), &values); err != nil {
		return metadata
	}
	delete(values, "media_url")
	encoded, err := json.Marshal(values)
	if err != nil {
		return metadata
	}
	return string(encoded)
}

// mediaViewURL signs the statement's media so the browser can load it. An
// unsigned or unconfigured storage just means no media is shown.
func mediaViewURL(app *appcontext.Context, e domain.Exercise) string {
	url := e.MediaURL()
	if url == "" || app == nil || app.ImageStorage == nil {
		return ""
	}
	signed, ok := app.ImageStorage.PresignGetURL(url, mediaLinkTTL)
	if !ok {
		return ""
	}
	return signed
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
