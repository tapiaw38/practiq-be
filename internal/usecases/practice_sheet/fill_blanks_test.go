package practicesheet

import "testing"

func TestBlanksAnswersMatch(t *testing.T) {
	t.Run("order of the keys does not change the verdict", func(t *testing.T) {
		// The student may fill blank 2 first; the JSON key order is arbitrary.
		if !blanksAnswersMatch(`{"2":"0","1":"pixel < 128"}`, `{"1":"pixel < 128","2":"0"}`) {
			t.Error("expected a match regardless of key order")
		}
	})

	t.Run("spacing is normalised on both sides", func(t *testing.T) {
		if !blanksAnswersMatch(`{"1":"pixel   <  128","2":" 0 "}`, `{"1":"pixel < 128","2":"0"}`) {
			t.Error("extra whitespace must not fail a correct answer")
		}
	})

	t.Run("values containing the old delimiters still work", func(t *testing.T) {
		// This is why the format is JSON: "||" is ordinary in code and ":" in
		// times or ratios. A delimited format broke on both.
		if !blanksAnswersMatch(`{"1":"a || b","2":"10:30"}`, `{"1":"a || b","2":"10:30"}`) {
			t.Error("options containing | or : must compare correctly")
		}
		if blanksAnswersMatch(`{"1":"a || b"}`, `{"1":"a"}`) {
			t.Error("a different value must not match")
		}
	})

	t.Run("works with prose, not just code", func(t *testing.T) {
		if !blanksAnswersMatch(`{"2":"1816","1":"Casa de Tucumán"}`, `{"1":"Casa de Tucumán","2":"1816"}`) {
			t.Error("expected a match for a history exercise")
		}
	})

	t.Run("a partial answer fails", func(t *testing.T) {
		// Filling one of two blanks is not a correct answer.
		if blanksAnswersMatch(`{"1":"Casa de Tucumán"}`, `{"1":"Casa de Tucumán","2":"1816"}`) {
			t.Error("a missing blank must fail")
		}
	})

	t.Run("an extra blank fails", func(t *testing.T) {
		if blanksAnswersMatch(`{"1":"a","2":"b"}`, `{"1":"a"}`) {
			t.Error("an answer with more blanks than expected must fail")
		}
	})

	t.Run("empty and malformed answers fail", func(t *testing.T) {
		for _, student := range []string{"", "   ", "no soy json", `{"1":""}`, `1:a|2:b`} {
			if blanksAnswersMatch(student, `{"1":"a","2":"b"}`) {
				t.Errorf("%q must not match", student)
			}
		}
	})

	t.Run("padded keys are the same blank", func(t *testing.T) {
		if !blanksAnswersMatch(`{"01":"a"}`, `{"1":"a"}`) {
			t.Error("blank 01 and blank 1 are the same blank")
		}
	})
}
