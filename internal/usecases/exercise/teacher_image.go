package exercise

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

func storeTeacherImage(ctx context.Context, app *appcontext.Context, ownerID, incoming string, previous domain.Exercise) string {
	values, ok := parseMetadata(incoming)
	if !ok {
		return incoming
	}

	key, value := findTeacherImage(values)
	switch {
	case key == "":

		if kept := previous.TeacherImage(); kept != "" {
			values["teacher_image"] = jsonString(kept)
			return encodeMetadata(values, incoming)
		}
		return incoming

	case value == "":

		return incoming

	case !strings.HasPrefix(value, "data:"):

		return incoming
	}

	if app.ImageStorage == nil {
		return incoming
	}
	uploaded, err := app.ImageStorage.UploadDataURI(ctx, "exercises", ownerID, value)
	if err != nil {

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
