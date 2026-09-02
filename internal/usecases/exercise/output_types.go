package exercise

import (
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

// mediaLinkTTL has to outlive solving the exercise without leaving a shareable
// link around for long. Reloading the sheet issues a fresh one.
const mediaLinkTTL = time.Hour

type (
	ExerciseData struct {
		ID            string `json:"id"`
		TopicID       string `json:"topic_id"`
		MaterialID    string `json:"material_id,omitempty"`
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
		// not here: it is fetched from the exercise's statement-image endpoint
		// by whoever needs to display it.
		HasTeacherImage bool   `json:"has_teacher_image,omitempty"`
		CreatedAt       string `json:"created_at"`
	}
)

func toExerciseData(app *appcontext.Context, e domain.Exercise) ExerciseData {
	return ExerciseData{
		ID:              e.ID,
		TopicID:         e.TopicID,
		MaterialID:      e.MaterialID,
		Type:            e.Type,
		Question:        e.Question,
		CorrectAnswer:   e.CorrectAnswer,
		Explanation:     e.Explanation,
		Difficulty:      e.Difficulty,
		Metadata:        e.MetadataWithoutTeacherImage(),
		MediaViewURL:    mediaViewURL(app, e),
		HasTeacherImage: e.TeacherImage() != "",
		CreatedAt:       e.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
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
