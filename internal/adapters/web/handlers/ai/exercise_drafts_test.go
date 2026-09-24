package ai

import (
	"strings"
	"testing"
)

func TestDraftMessageBodyWithoutFileSendsOnlyThePrompt(t *testing.T) {
	body, contentType, err := draftMessageBody(nil, "", false, "3", "7", "fracciones equivalentes", "")
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
	body, contentType, err := draftMessageBody([]byte("%PDF-1.7 fake"), "guia.pdf", true, "5", "5", "", "")
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
	body, contentType, err := draftMessageBody([]byte("\x89PNG fake"), "pizarron.png", false, "5", "5", "", "")
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

func TestDraftMessageBodyAsksForTheChosenTypeOnly(t *testing.T) {
	body, contentType, err := draftMessageBody(nil, "", false, "3", "4", "el agua", "fill_blanks")
	if err != nil {
		t.Fatalf("draftMessageBody() error = %v", err)
	}

	prompt := readMultipartValues(t, contentType, body)["content"]
	if !strings.Contains(prompt, `Generá únicamente ejercicios de tipo "fill_blanks"`) {
		t.Fatalf("prompt = %q, want it to pin the type", prompt)
	}
	if !strings.Contains(prompt, "{{1}}") {
		t.Fatalf("prompt = %q, want the worked example that shows the markers", prompt)
	}
}

func TestDraftMessageBodyWithoutATypeAsksForAMix(t *testing.T) {
	body, contentType, err := draftMessageBody(nil, "", false, "3", "4", "el agua", "")
	if err != nil {
		t.Fatalf("draftMessageBody() error = %v", err)
	}

	prompt := readMultipartValues(t, contentType, body)["content"]
	if !strings.Contains(prompt, "Variá los tipos") {
		t.Fatalf("prompt = %q, want the mixed-type instruction", prompt)
	}
	if strings.Contains(prompt, "únicamente ejercicios de tipo") {
		t.Fatalf("prompt = %q, want no type pinned when none was chosen", prompt)
	}
}

func TestDraftTypeRulePinsCanvas(t *testing.T) {
	if got := draftTypeRule("canvas"); !strings.Contains(got, `únicamente ejercicios de tipo "canvas"`) {
		t.Fatalf("draftTypeRule(\"canvas\") = %q, want canvas-only instruction", got)
	}
}

func TestValidateDraftShapeRejectsHandwritten(t *testing.T) {
	err := validateDraftShape(map[string]interface{}{
		"type":     "handwritten",
		"question": "Escribí esto a mano",
	}, "")
	if err == nil {
		t.Fatal("validateDraftShape() accepted handwritten draft")
	}
}

func TestValidateDraftShapeAcceptsManualCompatibleTypes(t *testing.T) {
	for _, typ := range []string{"open_text", "multiple_choice", "equation", "canvas", "attachment", "fill_blanks"} {
		if err := validateDraftShape(map[string]interface{}{"type": typ, "question": "Consigna"}, ""); err != nil {
			t.Fatalf("validateDraftShape(%q) error = %v", typ, err)
		}
	}
}
