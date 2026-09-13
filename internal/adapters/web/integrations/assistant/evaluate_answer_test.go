package assistant

import "testing"

func TestEvaluationRequiresExplicitVerdict(t *testing.T) {
	for _, raw := range []string{`{}`, `{"feedback":"uncertain"}`, `{"is_correct":null}`, `{"is_correct":"false"}`} {
		if _, err := parseEvaluationResponse(raw); err == nil {
			t.Errorf("accepted indeterminate evaluation: %s", raw)
		}
	}
	for _, verdict := range []string{"true", "false"} {
		result, err := parseEvaluationResponse(`{"is_correct":` + verdict + `,"feedback":"checked"}`)
		if err != nil || result.IsCorrect != (verdict == "true") {
			t.Fatalf("explicit %s: result=%+v err=%v", verdict, result, err)
		}
	}
}
