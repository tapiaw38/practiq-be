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

// TeacherImage is the statement the teacher drew by hand, stored in Metadata
// as a base64 data URL. Empty when the exercise has none.
//
// It is deliberately not part of the review list payload: a canvas is hundreds
// of KB of base64, and a hundred of them would make the queue unusable. The
// list reports whether one exists and it is fetched when actually opened.
func (e Exercise) TeacherImage() string {
	if e.Metadata == "" {
		return ""
	}
	var parsed map[string]any
	if err := json.Unmarshal([]byte(e.Metadata), &parsed); err != nil {
		return ""
	}
	for _, key := range teacherImageKeys {
		if value, ok := parsed[key].(string); ok && strings.HasPrefix(value, "data:image/") {
			return value
		}
	}
	return ""
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
