package practicesheet

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func TestMetadataWithoutMediaURLKeepsExerciseConfiguration(t *testing.T) {
	metadata := `{"media_url":"https://bucket/file/exercises/teacher/video.mp4","blanks":{"1":"5"},"accept":["audio"]}`

	got := metadataWithoutMediaURL(metadata)
	var values map[string]json.RawMessage
	if err := json.Unmarshal([]byte(got), &values); err != nil {
		t.Fatalf("returned invalid metadata: %v", err)
	}
	if _, ok := values["media_url"]; ok {
		t.Fatal("canonical media URL must not be returned in a practice sheet")
	}
	if string(values["blanks"]) != `{"1":"5"}` || string(values["accept"]) != `["audio"]` {
		t.Fatalf("non-media configuration was changed: %s", got)
	}
}

func TestToSheetDataHidesTeacherOnlyFieldsFromStudents(t *testing.T) {
	ps := domain.PracticeSheet{Exercises: []domain.PracticeSheetExercise{{
		Exercise: domain.Exercise{
			ID:            "exercise-1",
			CorrectAnswer: "42",
			Explanation:   "sumalos",
			Metadata:      `{"media_url":"https://bucket/file/exercises/teacher/image.png","options":["42"]}`,
			CreatedAt:     time.Now(),
		},
	}}}

	student := toSheetData(nil, ps, false).Exercises[0].Exercise
	if student.CorrectAnswer != "" || student.Explanation != "" {
		t.Fatal("student response must not include a solution or explanation")
	}
	if got := metadataWithoutMediaURL(student.Metadata); got != student.Metadata {
		t.Fatalf("student metadata still contains media URL: %s", student.Metadata)
	}
	encoded, err := json.Marshal(student)
	if err != nil {
		t.Fatalf("marshal student exercise: %v", err)
	}
	if strings.Contains(string(encoded), "correct_answer") || strings.Contains(string(encoded), "explanation") || strings.Contains(string(encoded), "media_url") {
		t.Fatalf("student JSON leaks teacher-only data: %s", encoded)
	}

	teacher := toSheetData(nil, ps, true).Exercises[0].Exercise
	if teacher.CorrectAnswer != "42" || teacher.Explanation != "sumalos" {
		t.Fatal("teacher response lost solution data")
	}
	if teacher.Metadata != ps.Exercises[0].Exercise.Metadata {
		t.Fatal("teacher response lost canonical metadata")
	}
}

func TestMetadataWithoutMediaURLLeavesInvalidMetadataUnchanged(t *testing.T) {
	if got := metadataWithoutMediaURL(`{invalid`); got != `{invalid` {
		t.Fatalf("got %q", got)
	}
}
