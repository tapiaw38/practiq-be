package practicesheet

import (
	"encoding/json"
	"strings"
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
		// HasTeacherImage says the statement was drawn by hand. The drawing is
		// fetched from the exercise's statement-image endpoint, never embedded:
		// one canvas outweighed the whole sheet around it.
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
		// ScheduledAt is empty when the sheet can be taken at any time.
		ScheduledAt string `json:"scheduled_at,omitempty"`
		// AvailableUntil closes the window; empty means it stays open.
		AvailableUntil string              `json:"available_until,omitempty"`
		CreatedBy      string              `json:"created_by"`
		CreatedAt      string              `json:"created_at"`
		Exercises      []SheetExerciseData `json:"exercises"`
	}

	ExerciseResultData struct {
		ExerciseID    string `json:"exercise_id"`
		IsCorrect     bool   `json:"is_correct"`
		StudentAnswer string `json:"student_answer"`
		CorrectAnswer string `json:"correct_answer"`
		AIFeedback    string `json:"ai_feedback,omitempty"`
		// NeedsTeacherReview marks an answer handed to the teacher. Only a level
		// test does that; on a practice it is always false.
		NeedsTeacherReview bool `json:"needs_teacher_review,omitempty"`
		// NotGraded marks an answer nobody could put a verdict on. It is left
		// out of the score rather than counted as wrong, so the student must be
		// told it was skipped instead of shown it as an error.
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
	data := sheetScalars(ps)
	data.Exercises = exercises
	return data
}

// toSheetSummary is the listing shape: the sheet, and the identity of the
// exercises on it, without their bodies.
//
// Every listing consumer only counts the exercises or collects their ids. The
// statement, its metadata and the signed media belong to the practice screen,
// which loads a single sheet by id. Sending them once per sheet made one course
// listing 46 KB, of which 93% was base64 images.
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
	}
	if ps.ScheduledAt != nil {
		data.ScheduledAt = ps.ScheduledAt.UTC().Format(timeFormat)
	}
	if ps.AvailableUntil != nil {
		data.AvailableUntil = ps.AvailableUntil.UTC().Format(timeFormat)
	}
	return data
}

// studentMetadata strips everything a student must not see before solving.
//
// Two things are removed:
//   - media_url, the stable bucket URL (MediaViewURL carries a signed one)
//   - the answer of every fill-blank, which is the solution itself
//
// Blank ids, the option pool and the layout stay: the student needs them to
// render the exercise, and the pool alone does not say which option goes where.
func studentMetadata(metadata string) string {
	if strings.TrimSpace(metadata) == "" {
		return ""
	}
	var values map[string]json.RawMessage
	if err := json.Unmarshal([]byte(metadata), &values); err != nil {
		// Metadata may hold fill-blank answers. A malformed value cannot be
		// redacted reliably, so fail closed instead of returning a secret.
		return "{}"
	}
	delete(values, "media_url")

	if raw, ok := values["blanks"]; ok {
		if redacted, err := redactBlankAnswers(raw); err == nil {
			values["blanks"] = redacted
		} else {
			// Unknown shape: drop it rather than risk shipping the answers.
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
