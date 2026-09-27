package practicesheet

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

const timeFormat = "2006-01-02T15:04:05Z"

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

		MediaViewURL string `json:"media_view_url,omitempty"`

		HasTeacherImage bool `json:"has_teacher_image,omitempty"`
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

		ScheduledAt string `json:"scheduled_at,omitempty"`

		AvailableUntil string `json:"available_until,omitempty"`

		MaxAttempts      *int `json:"max_attempts"`
		TimeLimitMinutes *int `json:"time_limit_minutes"`

		AttemptsUsed    int                 `json:"attempts_used,omitempty"`
		AttemptsAllowed int                 `json:"attempts_allowed,omitempty"`
		Deadline        string              `json:"deadline,omitempty"`
		CreatedBy       string              `json:"created_by"`
		CreatedAt       string              `json:"created_at"`
		Exercises       []SheetExerciseData `json:"exercises"`

		StreakDays int `json:"streak_days"`
	}

	ExerciseResultData struct {
		ExerciseID    string `json:"exercise_id"`
		IsCorrect     bool   `json:"is_correct"`
		StudentAnswer string `json:"student_answer"`
		CorrectAnswer string `json:"correct_answer"`
		AIFeedback    string `json:"ai_feedback,omitempty"`

		NeedsTeacherReview bool `json:"needs_teacher_review,omitempty"`

		NotGraded bool `json:"not_graded,omitempty"`
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
		StreakDays      int                  `json:"streak_days"`
		XPGained        int                  `json:"xp_gained"`
		CourseXP        int                  `json:"course_xp"`
		XPBreakdown     []XPBreakdownEntry   `json:"xp_breakdown,omitempty"`
		ExerciseResults []ExerciseResultData `json:"exercise_results"`
	}
)

func toSheetData(app *appcontext.Context, ps domain.PracticeSheet, includeTeacherData bool) PracticeSheetData {
	exercises := make([]SheetExerciseData, 0, len(ps.Exercises))
	for _, pse := range ps.Exercises {
		metadata := studentMetadata(pse.Exercise.MetadataWithoutTeacherImage())
		if includeTeacherData {
			metadata = pse.Exercise.MetadataWithoutTeacherImage()
		}
		exercise := ExerciseData{
			ID:              pse.Exercise.ID,
			TopicID:         pse.Exercise.TopicID,
			Type:            pse.Exercise.Type,
			Question:        pse.Exercise.Question,
			Difficulty:      pse.Exercise.Difficulty,
			Metadata:        metadata,
			MediaViewURL:    mediaViewURL(app, pse.Exercise),
			HasTeacherImage: pse.Exercise.TeacherImage() != "",
		}

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
	data := sheetScalars(ps)
	data.Exercises = exercises
	return data
}

func toSheetSummary(ps domain.PracticeSheet) PracticeSheetData {
	exercises := make([]SheetExerciseData, 0, len(ps.Exercises))
	for _, pse := range ps.Exercises {
		exercises = append(exercises, SheetExerciseData{
			ID:         pse.ID,
			OrderIndex: pse.OrderIndex,
			Exercise: ExerciseData{
				ID:      pse.Exercise.ID,
				TopicID: pse.Exercise.TopicID,
				Type:    pse.Exercise.Type,
			},
		})
	}
	data := sheetScalars(ps)
	data.Exercises = exercises
	return data
}

func sheetScalars(ps domain.PracticeSheet) PracticeSheetData {
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

		MaxAttempts:      ps.MaxAttempts,
		TimeLimitMinutes: ps.TimeLimitMinutes,
	}
	if ps.ScheduledAt != nil {
		data.ScheduledAt = ps.ScheduledAt.UTC().Format(timeFormat)
	}
	if ps.AvailableUntil != nil {
		data.AvailableUntil = ps.AvailableUntil.UTC().Format(timeFormat)
	}
	return data
}

func studentMetadata(metadata string) string {
	if strings.TrimSpace(metadata) == "" {
		return ""
	}
	var values map[string]json.RawMessage
	if err := json.Unmarshal([]byte(metadata), &values); err != nil {

		return "{}"
	}
	delete(values, "media_url")

	if raw, ok := values["blanks"]; ok {
		if redacted, err := redactBlankAnswers(raw); err == nil {
			values["blanks"] = redacted
		} else {

			delete(values, "blanks")
		}
	}

	encoded, err := json.Marshal(values)
	if err != nil {
		return "{}"
	}
	return string(encoded)
}

func redactBlankAnswers(raw json.RawMessage) (json.RawMessage, error) {
	var blanks []map[string]json.RawMessage
	if err := json.Unmarshal(raw, &blanks); err != nil {
		return nil, err
	}
	for _, blank := range blanks {
		delete(blank, "answer")
	}
	return json.Marshal(blanks)
}

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
