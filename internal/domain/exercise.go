package domain

import (
	"encoding/json"
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
