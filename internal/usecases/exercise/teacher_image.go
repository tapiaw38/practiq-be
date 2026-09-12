package exercise

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

// storeTeacherImage moves a handwritten statement out of the metadata and into
// the bucket, and keeps the one already stored when a save does not mention it.
//
// Both halves matter:
//
// The editor posts the drawing as a base64 data URL, which is how three of them
// came to weigh 25 KB inside a JSONB column and travel in every payload that
// touched the exercise. Uploading it here means the editor did not have to
// change and no client can reintroduce the inline form.
//
// Carrying the previous value over is what makes the change safe. Responses no
// longer include the drawing, so an editor that loads an exercise and saves it
// back sends metadata without it; without this, opening and saving an exercise
// would silently erase its statement.
func storeTeacherImage(ctx context.Context, app *appcontext.Context, ownerID, incoming string, previous domain.Exercise) string {
	values, ok := parseMetadata(incoming)
	if !ok {
		return incoming
	}

	key, value := findTeacherImage(values)
	switch {
	case key == "":
		// Not mentioned: keep whatever the exercise already had.
		if kept := previous.TeacherImage(); kept != "" {
			values["teacher_image"] = jsonString(kept)
			return encodeMetadata(values, incoming)
		}
		return incoming

	case value == "":
		// Mentioned as empty: the teacher cleared the drawing.
		return incoming

	case !strings.HasPrefix(value, "data:"):
		// Already a stored URL; nothing to upload.
		return incoming
	}

	if app.ImageStorage == nil {
		return incoming
	}
	uploaded, err := app.ImageStorage.UploadDataURI(ctx, "exercises", ownerID, value)
	if err != nil {
		// Falling back to the inline form keeps the teacher's work rather than
		// losing a drawing over a storage hiccup.
		log.Printf("[image_storage] exercise statement upload failed owner_id=%s err=%v", ownerID, err)
		return incoming
	}
	values[key] = jsonString(uploaded)
	return encodeMetadata(values, incoming)
}

func parseMetadata(metadata string) (map[string]json.RawMessage, bool) {
	if strings.TrimSpace(metadata) == "" {
		return nil, false
	}
	var values map[string]json.RawMessage
	if err := json.Unmarshal([]byte(metadata), &values); err != nil {
		return nil, false
	}
	return values, true
}

// findTeacherImage returns the key the drawing is under and its value. The key
// has been renamed over time, so which one is present says how old the payload
// is; writing back under the same key avoids leaving two of them behind.
func findTeacherImage(values map[string]json.RawMessage) (string, string) {
	for _, key := range []string{"teacher_image", "teacherImage", "image_data", "imageData"} {
		raw, ok := values[key]
		if !ok {
			continue
		}
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			return key, ""
		}
		return key, value
	}
	return "", ""
}

func encodeMetadata(values map[string]json.RawMessage, fallback string) string {
	encoded, err := json.Marshal(values)
	if err != nil {
		return fallback
	}
	return string(encoded)
}

func jsonString(value string) json.RawMessage {
	encoded, err := json.Marshal(value)
	if err != nil {
		return json.RawMessage(`""`)
	}
	return encoded
}
