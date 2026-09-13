package storage

import (
	"errors"
	"strings"
	"testing"
)

func TestClassifyContentType(t *testing.T) {
	cases := []struct {
		contentType string
		kind        FileKind
		ext         string
	}{
		{"application/pdf", FileKindPDF, ".pdf"},
		{"AUDIO/WEBM", FileKindAudio, ".webm"},
		// Browsers append codec parameters when recording.
		{"audio/webm;codecs=opus", FileKindAudio, ".webm"},
		{" image/jpeg ", FileKindImage, ".jpg"},
		{"application/vnd.openxmlformats-officedocument.wordprocessingml.document", FileKindDocument, ".docx"},
	}

	for _, tc := range cases {
		kind, ext, err := ClassifyContentType(tc.contentType)
		if err != nil {
			t.Errorf("%q: unexpected error: %v", tc.contentType, err)
			continue
		}
		if kind != tc.kind || ext != tc.ext {
			t.Errorf("%q: expected (%s, %s), got (%s, %s)", tc.contentType, tc.kind, tc.ext, kind, ext)
		}
	}
}

func TestResolveContentType(t *testing.T) {
	pdf := []byte("%PDF-1.4\n trailing bytes so sniffing has something to read")

	t.Run("a valid declared type is kept", func(t *testing.T) {
		contentType, kind, ext, err := ResolveContentType("application/pdf", pdf)
		if err != nil || contentType != "application/pdf" || kind != FileKindPDF || ext != ".pdf" {
			t.Errorf("got (%q, %s, %s, %v)", contentType, kind, ext, err)
		}
	})

	t.Run("an empty type falls back to sniffing", func(t *testing.T) {
		// Browsers sometimes send no type at all; the file must still land with
		// the type it really is, since that is what gets stored and reported.
		contentType, kind, _, err := ResolveContentType("", pdf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if kind != FileKindPDF {
			t.Errorf("expected the sniffed pdf kind, got %s (%q)", kind, contentType)
		}
	})

	t.Run("unsniffable garbage is rejected", func(t *testing.T) {
		if _, _, _, err := ResolveContentType("", []byte{0x00, 0x01, 0x02, 0x03}); !errors.Is(err, ErrUnsupportedFileType) {
			t.Errorf("expected rejection, got %v", err)
		}
	})
}

func TestBuildFileKeyIsNotFiledUnderImage(t *testing.T) {
	key := buildFileKey("attachments", "student-1", ".pdf")
	if strings.HasPrefix(key, "image/") {
		t.Errorf("non-image uploads must not live under image/: %s", key)
	}
	if !strings.HasPrefix(key, "file/attachments/student-1/") || !strings.HasSuffix(key, ".pdf") {
		t.Errorf("unexpected key layout: %s", key)
	}
	// Path parts are sanitised, so a crafted user id cannot escape the prefix.
	escaped := buildFileKey("attachments", "../../etc", ".pdf")
	if strings.Contains(escaped, "..") {
		t.Errorf("key must not contain traversal segments: %s", escaped)
	}
}

func TestClassifyContentTypeRejectsUnknown(t *testing.T) {
	// The whitelist is what keeps executables and scripts out of the bucket.
	for _, contentType := range []string{
		"application/x-msdownload",
		"application/octet-stream",
		"text/html",
		"",
	} {
		if _, _, err := ClassifyContentType(contentType); !errors.Is(err, ErrUnsupportedFileType) {
			t.Errorf("%q should be rejected, got err=%v", contentType, err)
		}
	}
}

// The stored content type decides the assistant channel and answers the
// exercise's accepted-format check, so a declaration the bytes contradict has
// to be rejected rather than trusted.
func TestResolveContentTypeRejectsMismatchedBytes(t *testing.T) {
	pdf := []byte("%PDF-1.4\n%âãÏÓ\nstartxref\n0\n%%EOF\n")

	t.Run("pdf bytes declared as audio", func(t *testing.T) {
		if _, _, _, err := ResolveContentType("audio/mpeg", pdf); err == nil {
			t.Fatal("expected the mismatch to be rejected")
		}
	})

	t.Run("pdf bytes declared as pdf", func(t *testing.T) {
		contentType, kind, _, err := ResolveContentType("application/pdf", pdf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if contentType != "application/pdf" || kind != FileKindPDF {
			t.Fatalf("got %s / %s", contentType, kind)
		}
	})

	t.Run("unsniffable bytes keep the declaration", func(t *testing.T) {
		// http.DetectContentType cannot identify these; rejecting here would
		// break every format Go does not know.
		opaque := []byte{0x1A, 0x45, 0xDF, 0xA3, 0x01, 0x02, 0x03, 0x04}
		contentType, kind, _, err := ResolveContentType("audio/webm", opaque)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if contentType != "audio/webm" || kind != FileKindAudio {
			t.Fatalf("got %s / %s", contentType, kind)
		}
	})
}
