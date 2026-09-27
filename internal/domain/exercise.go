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
	Metadata      string
	CreatedAt     time.Time
}

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

var teacherImageKeys = []string{"teacher_image", "teacherImage", "image_data", "imageData"}

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

func (e Exercise) MetadataWithoutTeacherImage() string {
	if strings.TrimSpace(e.Metadata) == "" {
		return e.Metadata
	}
	var values map[string]json.RawMessage
	if err := json.Unmarshal([]byte(e.Metadata), &values); err != nil {

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
