package invitecode_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/practiq-be/internal/platform/invitecode"
)

func TestNew(t *testing.T) {
	seen := make(map[string]bool, 500)

	for range 500 {
		code, err := invitecode.New()

		assert.NoError(t, err)
		assert.Len(t, code, invitecode.Length)
		assert.NotContains(t, code, "I")
		assert.NotContains(t, code, "L")
		assert.NotContains(t, code, "O")
		assert.NotContains(t, code, "U")
		assert.Equal(t, strings.ToUpper(code), code)
		assert.False(t, seen[code], "repeated code within 500 draws: %s", code)

		seen[code] = true
	}
}

func TestNormalize(t *testing.T) {
	tests := map[string]struct {
		raw      string
		expected string
	}{
		"already normalized":    {raw: "7K4P2QDX", expected: "7K4P2QDX"},
		"lowercase":             {raw: "7k4p2qdx", expected: "7K4P2QDX"},
		"with dash":             {raw: "7K4P-2QDX", expected: "7K4P2QDX"},
		"with spaces":           {raw: "  7K4P 2QDX ", expected: "7K4P2QDX"},
		"drops foreign letters": {raw: "7K4P-2QDX!io", expected: "7K4P2QDX"},
		"empty":                 {raw: "", expected: ""},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, invitecode.Normalize(tc.raw))
		})
	}
}

func TestFormat(t *testing.T) {
	tests := map[string]struct {
		code     string
		expected string
	}{
		"adds the dash in the middle":   {code: "7K4P2QDX", expected: "7K4P-2QDX"},
		"leaves unexpected input as is": {code: "SHORT", expected: "SHORT"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, invitecode.Format(tc.code))
		})
	}
}

// Un código tiene que sobrevivir el viaje por la pantalla: se muestra con guion
// y el alumno lo pega tal cual.
func TestFormatThenNormalizeRoundTrip(t *testing.T) {
	code, err := invitecode.New()

	assert.NoError(t, err)
	assert.Equal(t, code, invitecode.Normalize(invitecode.Format(code)))
}
