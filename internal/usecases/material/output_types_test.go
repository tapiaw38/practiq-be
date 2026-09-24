package material

import (
	"strings"
	"testing"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func TestToMaterialPreviewCutsWithoutBreakingAccents(t *testing.T) {
	// The repository hands back previewChars+1 characters, which is how a text
	// that merely ends at the limit is told apart from one that continues.
	cases := []struct {
		name          string
		text          string
		wantLen       int
		wantTruncated bool
	}{
		{
			name:          "a short text is served whole",
			text:          "Sustantivos y su clasificación",
			wantLen:       30,
			wantTruncated: false,
		},
		{
			name:          "a text ending exactly at the limit is not marked",
			text:          strings.Repeat("a", previewChars),
			wantLen:       previewChars,
			wantTruncated: false,
		},
		{
			name:          "one character past the limit is marked",
			text:          strings.Repeat("a", previewChars+1),
			wantLen:       previewChars,
			wantTruncated: true,
		},
		{
			// Cutting bytes instead of runes would leave half an "ó" behind and
			// the response would not be valid UTF-8.
			name:          "accents are not split in half",
			text:          strings.Repeat("ó", previewChars+1),
			wantLen:       previewChars,
			wantTruncated: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := toMaterialPreview(domain.Material{ExtractedText: tc.text})
			if n := len([]rune(got.ExtractedText)); n != tc.wantLen {
				t.Fatalf("got %d characters, want %d", n, tc.wantLen)
			}
			if got.ExtractedTextTruncated != tc.wantTruncated {
				t.Fatalf("truncated = %v, want %v", got.ExtractedTextTruncated, tc.wantTruncated)
			}
			if !strings.HasPrefix(tc.text, got.ExtractedText) {
				t.Fatal("the preview is not the beginning of the text")
			}
		})
	}
}
