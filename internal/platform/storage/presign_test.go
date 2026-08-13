package storage

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/tapiaw38/practiq-be/internal/platform/config"
)

func testStorage() *S3ImageStorage {
	return &S3ImageStorage{cfg: config.S3Config{
		AWSRegion:          "us-east-1",
		AWSBucket:          "practiq",
		AWSAccessKeyID:     "AKIAEXAMPLE",
		AWSSecretAccessKey: "secret",
	}}
}

func TestPresignGetURLSignsStoredObject(t *testing.T) {
	s := testStorage()
	stored := "https://practiq.s3.us-east-1.amazonaws.com/file/materials/user1/abc.pdf"

	signed, ok := s.PresignGetURL(stored, 15*time.Minute)
	if !ok {
		t.Fatalf("expected the stored object to be presigned")
	}

	u, err := url.Parse(signed)
	if err != nil {
		t.Fatalf("presigned URL does not parse: %v", err)
	}
	if u.Path != "/file/materials/user1/abc.pdf" {
		t.Errorf("path changed: got %q", u.Path)
	}

	q := u.Query()
	for _, key := range []string{
		"X-Amz-Algorithm", "X-Amz-Credential", "X-Amz-Date",
		"X-Amz-Expires", "X-Amz-SignedHeaders", "X-Amz-Signature",
	} {
		if q.Get(key) == "" {
			t.Errorf("missing %s", key)
		}
	}
	if got := q.Get("X-Amz-Expires"); got != "900" {
		t.Errorf("expires: got %q, want 900", got)
	}
	if got := q.Get("X-Amz-Credential"); !strings.HasPrefix(got, "AKIAEXAMPLE/") ||
		!strings.HasSuffix(got, "/us-east-1/s3/aws4_request") {
		t.Errorf("credential scope: got %q", got)
	}
	if got := len(q.Get("X-Amz-Signature")); got != 64 {
		t.Errorf("signature length: got %d, want 64 hex chars", got)
	}
}

// A material can hold a link we never stored; signing it would corrupt it.
func TestPresignGetURLLeavesForeignURLsAlone(t *testing.T) {
	s := testStorage()
	foreign := "https://www.youtube.com/watch?v=abc"

	got, ok := s.PresignGetURL(foreign, time.Hour)
	if ok || got != foreign {
		t.Fatalf("got (%q, %v), want the URL unchanged and false", got, ok)
	}
}

func TestPresignGetURLDefaultsTTL(t *testing.T) {
	s := testStorage()
	signed, ok := s.PresignGetURL("https://practiq.s3.us-east-1.amazonaws.com/file/a.pdf", 0)
	if !ok {
		t.Fatal("expected a signed URL")
	}
	u, _ := url.Parse(signed)
	if got := u.Query().Get("X-Amz-Expires"); got != "900" {
		t.Errorf("default expires: got %q, want 900", got)
	}
}
