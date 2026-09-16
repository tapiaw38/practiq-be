package upload

import (
	"testing"

	"github.com/tapiaw38/practiq-be/internal/platform/storage"
)

// A statement's material used to be images and audio only, so a teacher whose
// exercise came on a PDF had to screenshot it. The picker was widened first
// and uploads kept answering 415, which is the case these pin down.
func TestAllowedAsStatementMaterial(t *testing.T) {
	for _, tc := range []struct {
		kind storage.FileKind
		want bool
	}{
		{storage.FileKindImage, true},
		{storage.FileKindAudio, true},
		{storage.FileKindPDF, true},
		{storage.FileKindDocument, true},
		// Nothing renders a video beside the statement and no assistant
		// channel reads one.
		{storage.FileKindVideo, false},
	} {
		if got := allowedAsStatementMaterial(tc.kind); got != tc.want {
			t.Errorf("allowedAsStatementMaterial(%q) = %v, want %v", tc.kind, got, tc.want)
		}
	}
}

// The rule runs on the kind the bytes resolve to, not on what the client
// declared, so a renamed binary cannot walk in behind a PDF's content type.
func TestStatementMaterialFollowsTheResolvedKind(t *testing.T) {
	pdf := []byte("%PDF-1.4\n and enough trailing bytes for sniffing to work")

	_, kind, _, err := storage.ResolveContentType("application/pdf", pdf)
	if err != nil {
		t.Fatalf("ResolveContentType() error = %v", err)
	}
	if !allowedAsStatementMaterial(kind) {
		t.Fatalf("a PDF resolved to %q, which a statement refuses", kind)
	}

	if _, _, _, err := storage.ResolveContentType("application/pdf", []byte("\x89PNG\r\n\x1a\n fake png bytes")); err == nil {
		t.Fatal("PNG bytes declared as application/pdf were accepted")
	}
}
