package notebook

import "testing"

// AddPage uploads the teacher's page image and stores the URL it gets back, so
// ContentData is normally an https link. isLikelyImageData only recognised
// base64, and a URL is not base64-like, so the link reached the model verbatim
// as "respuesta correcta esperada: https://….png" and every submission came
// back UNREADABLE.
func TestExpectedAnswerHidesAnUploadedPageImage(t *testing.T) {
	uploaded := "https://practiq-images-264914792937-sa-east-1-an.s3.sa-east-1.amazonaws.com/image/notebook/xQDQCpOzZzLkBQPqQIvIFUSRRGquCv/85154628b0a1f1c2ab2e75d49df0e6fb.png"
	if got := normalizeNotebookExpectedAnswer(uploaded); got != "[imagen del docente]" {
		t.Fatalf("uploaded page image leaked into the prompt: %q", got)
	}
}

func TestExpectedAnswerKeepsRealText(t *testing.T) {
	for _, answer := range []string{
		"4",
		"x = 5",
		"Resolve las sumas de la pagina",
		"Ver https://es.wikipedia.org/wiki/Suma para repasar",
	} {
		if got := normalizeNotebookExpectedAnswer(answer); got != answer {
			t.Fatalf("expected %q to pass through, got %q", answer, got)
		}
	}
}

func TestIsImageURL(t *testing.T) {
	images := []string{
		"https://bucket.s3.amazonaws.com/a/b.png",
		"http://bucket.s3.amazonaws.com/a/b.JPG",
		"https://bucket.s3.amazonaws.com/a/b.jpeg?X-Amz-Signature=abc",
		"https://bucket.s3.amazonaws.com/a/b.webp#frag",
	}
	for _, value := range images {
		if !isImageURL(value) {
			t.Fatalf("expected %q to be treated as an image URL", value)
		}
	}

	notImages := []string{
		"",
		"4",
		"data:image/png;base64,iVBORw0KGgo",
		"https://example.com/page",
		"https://example.com/notes.pdf",
		"mira esto https://example.com/a.png y resolve",
	}
	for _, value := range notImages {
		if isImageURL(value) {
			t.Fatalf("expected %q not to be treated as an image URL", value)
		}
	}
}
