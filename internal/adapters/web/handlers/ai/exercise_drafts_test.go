package ai

import (
	"strings"
	"testing"
)

func TestDraftMessageBodyWithoutFileSendsOnlyThePrompt(t *testing.T) {
	body, contentType, err := draftMessageBody(nil, "", false, "3", "7", "fracciones equivalentes")
	if err != nil {
		t.Fatalf("draftMessageBody() error = %v", err)
	}

	values := readMultipartValues(t, contentType, body)
	prompt := values["content"]
	if !strings.Contains(prompt, "sobre el tema indicado abajo") {
		t.Fatalf("prompt = %q, want it to describe a topic rather than an attachment", prompt)
	}
	if strings.Contains(prompt, "archivo adjunto") {
		t.Fatalf("prompt = %q, want no mention of an attachment when none was sent", prompt)
	}
	if !strings.Contains(prompt, "fracciones equivalentes") {
		t.Fatalf("prompt = %q, want the teacher's instruction carried through", prompt)
	}
	if _, sent := values["document_content"]; sent {
		t.Fatal("a document part was sent even though there is no file")
	}
	if _, sent := values["image_content"]; sent {
		t.Fatal("an image part was sent even though there is no file")
	}
}

func TestDraftMessageBodyWithFileKeepsTheAttachment(t *testing.T) {
	body, contentType, err := draftMessageBody([]byte("%PDF-1.7 fake"), "guia.pdf", true, "5", "5", "")
	if err != nil {
		t.Fatalf("draftMessageBody() error = %v", err)
	}

	values := readMultipartValues(t, contentType, body)
	if !strings.Contains(values["content"], "archivo adjunto") {
		t.Fatalf("prompt = %q, want it to point at the attachment", values["content"])
	}
	if values["document_content"] != "%PDF-1.7 fake" {
		t.Fatalf("document_content = %q, want the uploaded bytes", values["document_content"])
	}
}

func TestDraftMessageBodySendsAnImagePartForImages(t *testing.T) {
	body, contentType, err := draftMessageBody([]byte("\x89PNG fake"), "pizarron.png", false, "5", "5", "")
	if err != nil {
		t.Fatalf("draftMessageBody() error = %v", err)
	}

	values := readMultipartValues(t, contentType, body)
	if values["image_content"] != "\x89PNG fake" {
		t.Fatalf("image_content = %q, want the uploaded bytes", values["image_content"])
	}
	if _, sent := values["document_content"]; sent {
		t.Fatal("an image was sent as a document")
	}
}
