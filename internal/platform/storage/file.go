package storage

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
)

// Storage is an alias for the historical ImageStorage name, which now also
// carries audio and documents.
type Storage = ImageStorage

// MaxUploadBytes caps a single upload. Audio recordings and scanned PDFs are
// the large cases; anything past this is likely a mistake.
//
// 50 MiB accommodates scans and longer audio without allowing unbounded
// request bodies. Larger media should move to multipart uploads.
const MaxUploadBytes = 50 << 20 // 50 MiB

// FileKind groups accepted content types into the buckets the product talks
// about, so exercises can say "accepts audio" instead of listing MIME types.
type FileKind string

const (
	FileKindAudio    FileKind = "audio"
	FileKindImage    FileKind = "image"
	FileKindPDF      FileKind = "pdf"
	FileKindDocument FileKind = "doc"
	FileKindVideo    FileKind = "video"
)

var ErrUnsupportedFileType = errors.New("unsupported file type")

// acceptedTypes is a whitelist: an unknown content type is rejected rather
// than stored, so the bucket never receives arbitrary binaries.
var acceptedTypes = map[string]struct {
	kind FileKind
	ext  string
}{
	"audio/webm":         {FileKindAudio, ".webm"},
	"audio/ogg":          {FileKindAudio, ".ogg"},
	"audio/mpeg":         {FileKindAudio, ".mp3"},
	"audio/mp4":          {FileKindAudio, ".m4a"},
	"audio/wav":          {FileKindAudio, ".wav"},
	"audio/x-wav":        {FileKindAudio, ".wav"},
	"audio/wave":         {FileKindAudio, ".wav"},
	"image/png":          {FileKindImage, ".png"},
	"image/jpeg":         {FileKindImage, ".jpg"},
	"image/jpg":          {FileKindImage, ".jpg"},
	"image/webp":         {FileKindImage, ".webp"},
	"image/gif":          {FileKindImage, ".gif"},
	"video/mp4":          {FileKindVideo, ".mp4"},
	"video/webm":         {FileKindVideo, ".webm"},
	"video/ogg":          {FileKindVideo, ".ogv"},
	"video/quicktime":    {FileKindVideo, ".mov"},
	"application/pdf":    {FileKindPDF, ".pdf"},
	"application/msword": {FileKindDocument, ".doc"},
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {FileKindDocument, ".docx"},
	"text/plain": {FileKindDocument, ".txt"},
	"application/vnd.oasis.opendocument.text": {FileKindDocument, ".odt"},
}

// ClassifyContentType reports the bucket a content type belongs to. Parameters
// such as "audio/webm;codecs=opus" are stripped before matching.
func ClassifyContentType(contentType string) (FileKind, string, error) {
	base := strings.ToLower(strings.TrimSpace(contentType))
	if index := strings.Index(base, ";"); index >= 0 {
		base = strings.TrimSpace(base[:index])
	}
	entry, ok := acceptedTypes[base]
	if !ok {
		return "", "", fmt.Errorf("%w: %s", ErrUnsupportedFileType, contentType)
	}
	return entry.kind, entry.ext, nil
}

// ResolveContentType returns the content type an upload will actually be
// stored with: the declared one when the whitelist accepts it *and* the bytes
// do not say otherwise, else the sniffed one.
//
// Checking the declaration alone used to be enough, because the type only
// decided how the object was served. It no longer is: the stored type now
// picks the assistant channel an attachment is sent through and answers the
// exercise's accepted-format check, so PDF bytes declared as audio/mpeg would
// slip past an audio-only exercise and be handed to the voice channel.
//
// http.DetectContentType only recognizes some formats; when it cannot tell
// (application/octet-stream and friends) the declaration stands, since
// rejecting there would break every format Go cannot sniff.
func ResolveContentType(contentType string, body []byte) (string, FileKind, string, error) {
	if kind, ext, err := ClassifyContentType(contentType); err == nil {
		if sniffedKind, _, sniffErr := ClassifyContentType(http.DetectContentType(body)); sniffErr == nil && conflictingKinds(kind, sniffedKind) {
			return "", "", "", fmt.Errorf("%w: declared %s but the file is %s", ErrUnsupportedFileType, contentType, sniffedKind)
		}
		return contentType, kind, ext, nil
	}
	// Some browsers send an empty or generic type; sniff the bytes before
	// rejecting the upload.
	sniffed := http.DetectContentType(body)
	kind, ext, err := ClassifyContentType(sniffed)
	if err != nil {
		return "", "", "", err
	}
	return sniffed, kind, ext, nil
}

// conflictingKinds reports a declaration the bytes contradict.
//
// Audio and video are never treated as a conflict: WebM and Ogg are the same
// container either way, so a voice recording is sniffed as video/webm and
// rejecting it would break the student's audio answers. Everything else that
// Go can identify has to agree.
func conflictingKinds(declared, sniffed FileKind) bool {
	if declared == sniffed {
		return false
	}
	mediaContainer := func(k FileKind) bool {
		return k == FileKindAudio || k == FileKindVideo
	}
	return !(mediaContainer(declared) && mediaContainer(sniffed))
}

// UploadFile stores a file and returns its URL. contentType is trusted only
// after passing the whitelist; when the client sends nothing usable the bytes
// are sniffed instead.
func (s *S3ImageStorage) UploadFile(ctx context.Context, folder, userID, filename, contentType string, body []byte) (string, error) {
	if len(body) == 0 {
		return "", errors.New("empty file")
	}
	if len(body) > MaxUploadBytes {
		return "", fmt.Errorf("file is larger than %d MiB", MaxUploadBytes>>20)
	}

	contentType, kind, ext, err := ResolveContentType(contentType, body)
	if err != nil {
		return "", err
	}
	if folder == "exercises" && kind != FileKindImage && kind != FileKindAudio {
		return "", fmt.Errorf("%w: exercise media must be image or audio", ErrUnsupportedFileType)
	}

	// The extension always comes from the verified content type, never from the
	// client-supplied filename.
	key := buildFileKey(folder, userID, ext)
	if err := s.putObject(ctx, key, contentType, body); err != nil {
		return "", err
	}
	return s.objectURL(key), nil
}

// buildFileKey mirrors buildImageKey but under a "file" prefix, so audio and
// documents do not end up filed under "image/". Existing objects keep working:
// lookups parse the full URL, not the prefix.
func buildFileKey(folder, userID, ext string) string {
	return path.Join("file", cleanPathPart(folder), cleanPathPart(userID), randomHex(16)+ext)
}

// UploadFile on the noop storage keeps local development working without S3
// credentials; there is nowhere to put the bytes, so it reports that.
func (NoopImageStorage) UploadFile(ctx context.Context, folder, userID, filename, contentType string, body []byte) (string, error) {
	return "", errors.New("file storage is not configured")
}

// FetchFile reads back a stored object, so an uploaded answer can be forwarded
// to the assistant for grading.
func (s *S3ImageStorage) FetchFile(ctx context.Context, url string) ([]byte, string, error) {
	key, ok := s.keyFromValue(url)
	if !ok {
		return nil, "", fmt.Errorf("%q is not a stored object", url)
	}
	body, contentType, err := s.getObject(ctx, key)
	if err != nil {
		return nil, "", err
	}
	if contentType == "" {
		contentType = http.DetectContentType(body)
	}
	return body, contentType, nil
}

func (NoopImageStorage) FetchFile(ctx context.Context, url string) ([]byte, string, error) {
	return nil, "", errors.New("file storage is not configured")
}
