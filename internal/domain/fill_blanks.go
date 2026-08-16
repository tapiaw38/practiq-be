package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Fill-in-the-blanks exercises: a statement in plain text with {{n}} markers,
// plus the answers and options in metadata. Subject-agnostic by design — the
// same type serves prose, dates, formulas or code.

var blankMarkerPattern = regexp.MustCompile(`\{\{\s*(\d+)\s*\}\}`)

var (
	ErrNoBlanks              = errors.New("the statement has no {{n}} blanks")
	ErrDuplicateMarker       = errors.New("a blank number is repeated in the statement")
	ErrDuplicateAnswerID     = errors.New("the answer repeats a blank number")
	ErrBlanksMismatch        = errors.New("the configured blanks do not match the statement")
	ErrEmptyBlankAnswer      = errors.New("a blank has no answer")
	ErrMissingOptions        = errors.New("the options do not cover every answer")
	ErrGradingAnswerMismatch = errors.New("the grading answer does not match the blanks in metadata")
)

// BlankIDsInStatement returns the blank numbers in the order they appear.
// A repeated marker is rejected: two blanks sharing a number would fill each
// other in the UI and read as one answer, which is never what the teacher meant.
func BlankIDsInStatement(statement string) ([]int, error) {
	matches := blankMarkerPattern.FindAllStringSubmatch(statement, -1)
	ids := make([]int, 0, len(matches))
	seen := map[int]bool{}

	for _, match := range matches {
		id, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		if seen[id] {
			return nil, fmt.Errorf("%w: {{%d}}", ErrDuplicateMarker, id)
		}
		seen[id] = true
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil, ErrNoBlanks
	}
	return ids, nil
}

// BlanksAnswer maps a blank number to what was placed in it.
type BlanksAnswer map[string]string

// ParseBlanksAnswer reads the JSON answer format, {"1":"a","2":"b"}. Values are
// trimmed and their inner whitespace collapsed so formatting cannot decide a
// grade.
//
// Keys are normalised ("01" and "1" are the same blank), and a collision
// between two keys that normalise to the same blank is an error rather than a
// coin flip: Go map iteration is random, so picking a winner would make the
// same answer grade differently on each attempt.
func ParseBlanksAnswer(value string) (BlanksAnswer, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	raw := map[string]string{}
	if err := json.Unmarshal([]byte(value), &raw); err != nil {
		return nil, err
	}

	parsed := make(BlanksAnswer, len(raw))
	// Keep this separate from parsed: an empty duplicate must still be rejected.
	// Otherwise map iteration could discard it before seeing its equivalent key.
	seenIDs := make(map[string]struct{}, len(raw))
	for id, answer := range raw {
		key := strings.TrimSpace(id)
		if number, err := strconv.Atoi(key); err == nil {
			key = strconv.Itoa(number)
		}
		if _, exists := seenIDs[key]; exists {
			return nil, fmt.Errorf("%w: %q", ErrDuplicateAnswerID, key)
		}
		seenIDs[key] = struct{}{}
		answer = NormalizeBlankAnswer(answer)
		if answer == "" {
			continue
		}
		parsed[key] = answer
	}
	if len(parsed) == 0 {
		return nil, nil
	}
	return parsed, nil
}

// NormalizeBlankAnswer collapses whitespace so spacing never decides a grade.
func NormalizeBlankAnswer(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

// BlanksAnswersMatch reports whether the student placed every blank as
// expected. A missing, extra or malformed blank fails.
func BlanksAnswersMatch(student, expected string) bool {
	got, err := ParseBlanksAnswer(student)
	if err != nil || got == nil {
		return false
	}
	want, err := ParseBlanksAnswer(expected)
	if err != nil || want == nil || len(got) != len(want) {
		return false
	}
	for id, answer := range want {
		if !strings.EqualFold(got[id], answer) {
			return false
		}
	}
	return true
}

// FillBlanksMetadata is the exercise configuration stored as JSON.
type FillBlanksMetadata struct {
	Blanks []struct {
		ID     int    `json:"id"`
		Answer string `json:"answer"`
	} `json:"blanks"`
	Options []string `json:"options"`
	Layout  string   `json:"layout"`
}

// ValidateFillBlanksGradingAnswer rejects a correct_answer that disagrees with
// the blanks in metadata. Grading compares the student's JSON only against
// correct_answer, so a mismatch here marks every right selection as wrong and
// nothing else in the system notices.
func ValidateFillBlanksGradingAnswer(correctAnswer, metadata string) error {
	var config FillBlanksMetadata
	if err := json.Unmarshal([]byte(strings.TrimSpace(metadata)), &config); err != nil {
		return fmt.Errorf("metadata is not valid JSON: %w", err)
	}

	placements := make(map[int]string, len(config.Blanks))
	for _, blank := range config.Blanks {
		placements[blank.ID] = NormalizeBlankAnswer(blank.Answer)
	}
	expected, err := json.Marshal(placements)
	if err != nil {
		return err
	}

	if !BlanksAnswersMatch(correctAnswer, string(expected)) {
		return ErrGradingAnswerMismatch
	}
	return nil
}

// ValidateFillBlanksExercise rejects an exercise a student could not solve.
// Enforced here rather than only in the UI, since the API is also reachable
// from the MCP server and from scripts.
func ValidateFillBlanksExercise(statement, metadata string) error {
	ids, err := BlankIDsInStatement(statement)
	if err != nil {
		return err
	}

	var config FillBlanksMetadata
	if err := json.Unmarshal([]byte(strings.TrimSpace(metadata)), &config); err != nil {
		return fmt.Errorf("metadata is not valid JSON: %w", err)
	}

	answers := make(map[int]string, len(config.Blanks))
	for _, blank := range config.Blanks {
		if _, exists := answers[blank.ID]; exists {
			return fmt.Errorf("%w: blank %d is configured twice", ErrBlanksMismatch, blank.ID)
		}
		answer := NormalizeBlankAnswer(blank.Answer)
		if answer == "" {
			return fmt.Errorf("%w: blank %d", ErrEmptyBlankAnswer, blank.ID)
		}
		answers[blank.ID] = answer
	}

	if len(answers) != len(ids) {
		return fmt.Errorf("%w: the statement has %d blanks and %d are configured", ErrBlanksMismatch, len(ids), len(answers))
	}
	for _, id := range ids {
		if _, ok := answers[id]; !ok {
			return fmt.Errorf("%w: blank %d has no answer", ErrBlanksMismatch, id)
		}
	}

	// Every answer needs its own block in the pool, counting repeats: two blanks
	// sharing an answer need two blocks or the student can only fill one.
	available := make(map[string]int, len(config.Options))
	for _, option := range config.Options {
		available[strings.ToLower(NormalizeBlankAnswer(option))]++
	}
	for id, answer := range answers {
		key := strings.ToLower(answer)
		if available[key] <= 0 {
			return fmt.Errorf("%w: no option available for blank %d", ErrMissingOptions, id)
		}
		available[key]--
	}
	return nil
}
