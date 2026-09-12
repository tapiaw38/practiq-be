package domain

import (
	"encoding/json"
	"strings"
	"time"
)

type Exercise struct {
	ID            string
	TopicID       string
	MaterialID    string
	Type          string
	Question      string
	CorrectAnswer string
	Explanation   string
	Difficulty    int
	Metadata      string // JSON string
	CreatedAt     time.Time
}

// AcceptedAttachmentKinds lists the file families the teacher allows as an
// answer, empty when any supported format goes. Stored in Metadata next to
// media_url; see MediaURL for why it lives there.
func (e Exercise) AcceptedAttachmentKinds() []string {
	if e.Metadata == "" {
		return nil
	}
	var parsed struct {
		Accept []string `json:"accept"`
	}
	if err := json.Unmarshal([]byte(e.Metadata), &parsed); err != nil {
		return nil
	}
	return parsed.Accept
}

// teacherImageKeys mirrors the frontend: the handwritten prompt has been
// stored under several names over time.
var teacherImageKeys = []string{"teacher_image", "teacherImage", "image_data", "imageData"}

// TeacherImage is the statement the teacher drew by hand, or empty when the
// exercise has none.
//
// Two forms are returned. Exercises saved from the current editor hold a bucket
// URL, like every other upload in the product. Older ones hold the drawing
// inline as a base64 data URL, which is why a single canvas could weigh 9 KB of
// JSON; those are still read here so they keep working, and turn into a URL the
// next time the exercise is saved.
//
// No response embeds this value. A canvas is orders of magnitude larger than
// the exercise around it, so every payload carries a flag and the image is
// fetched on its own when something actually draws it.
func (e Exercise) TeacherImage() string {
	if e.Metadata == "" {
		return ""
	}
	var parsed map[string]any
	if err := json.Unmarshal([]byte(e.Metadata), &parsed); err != nil {
		return ""
	}
	for _, key := range teacherImageKeys {
		value, ok := parsed[key].(string)
		if !ok {
			continue
		}
		if strings.HasPrefix(value, "data:image/") || strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
			return value
		}
	}
	return ""
}

// MetadataWithoutTeacherImage is Metadata with the drawing taken out, which is
// the only form that belongs in a response.
//
// The flag callers report alongside it comes from TeacherImage, so a stripped
// payload still says an image exists; it just does not carry it.
func (e Exercise) MetadataWithoutTeacherImage() string {
	if strings.TrimSpace(e.Metadata) == "" {
		return e.Metadata
	}
	var values map[string]json.RawMessage
	if err := json.Unmarshal([]byte(e.Metadata), &values); err != nil {
		// Unparseable metadata may be a bare data URL, which earlier versions
		// stored. Returning it would ship the very thing this strips.
		return "{}"
	}
	found := false
	for _, key := range teacherImageKeys {
		if _, ok := values[key]; ok {
			delete(values, key)
			found = true
		}
	}
	if !found {
		return e.Metadata
	}
	encoded, err := json.Marshal(values)
	if err != nil {
		return "{}"
	}
	return string(encoded)
}

// MediaURL is the image/audio/video the teacher attached to the statement, or
// empty when there is none. It lives inside Metadata (like teacher_image and
// accept already do) so attaching media needed no schema change.
//
// The value is the canonical storage URL: the bucket is private, so it has to
// be presigned before a browser can load it.
func (e Exercise) MediaURL() string {
	if e.Metadata == "" {
		return ""
	}
	var parsed struct {
		MediaURL string `json:"media_url"`
	}
	if err := json.Unmarshal([]byte(e.Metadata), &parsed); err != nil {
		return ""
	}
	return parsed.MediaURL
}
