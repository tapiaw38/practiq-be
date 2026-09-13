package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"strings"
	"testing"
)

func TestEnrichTutorMessageAddsInstructionAndPreservesParts(t *testing.T) {
	body, contentType := multipartBody(t, map[string]string{
		"content": "Necesito ayuda con fracciones",
		"context": "Ejercicio 2",
	})

	rewrittenType, rewrittenBody, err := enrichTutorMessage(contentType, body)
	if err != nil {
		t.Fatalf("enrichTutorMessage() error = %v", err)
	}

	values := readMultipartValues(t, rewrittenType, rewrittenBody)
	if got := values["context"]; got != "Ejercicio 2" {
		t.Fatalf("context = %q, want preserved value", got)
	}
	if got := values["content"]; got != tutorInstruction+"\nNecesito ayuda con fracciones" {
		t.Fatalf("content = %q, want enriched message", got)
	}
}

func TestEnrichTutorMessageDoesNotDuplicateLegacyInstruction(t *testing.T) {
	legacy := "POLITICA OBLIGATORIA:\nMensaje del alumno: hola"
	body, contentType := multipartBody(t, map[string]string{"content": legacy})

	rewrittenType, rewrittenBody, err := enrichTutorMessage(contentType, body)
	if err != nil {
		t.Fatalf("enrichTutorMessage() error = %v", err)
	}

	if got := readMultipartValues(t, rewrittenType, rewrittenBody)["content"]; got != legacy {
		t.Fatalf("content = %q, want legacy content unchanged", got)
	}
}

func TestEnrichTutorMessageRejectsInvalidMultipart(t *testing.T) {
	if _, _, err := enrichTutorMessage("multipart/form-data", []byte("invalid")); err == nil {
		t.Fatal("enrichTutorMessage() error = nil, want error")
	}
}

func TestTextMessageToTutorMultipart(t *testing.T) {
	body, err := json.Marshal(textMessageInput{
		Content: "¿Cómo empiezo?",
		Context: "Tema: fracciones",
	})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	contentType, multipartBody, err := textMessageToTutorMultipart("application/json", body)
	if err != nil {
		t.Fatalf("textMessageToTutorMultipart() error = %v", err)
	}
	values := readMultipartValues(t, contentType, multipartBody)
	if got := values["context"]; got != "Tema: fracciones" {
		t.Fatalf("context = %q, want preserved value", got)
	}
	if got := values["content"]; got != tutorInstruction+"\n¿Cómo empiezo?" {
		t.Fatalf("content = %q, want enriched message", got)
	}
}

func multipartBody(t *testing.T, values map[string]string) ([]byte, string) {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for name, value := range values {
		if err := writer.WriteField(name, value); err != nil {
			t.Fatalf("WriteField() error = %v", err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	return body.Bytes(), writer.FormDataContentType()
}

func readMultipartValues(t *testing.T, contentType string, body []byte) map[string]string {
	t.Helper()
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		t.Fatalf("ParseMediaType() error = %v", err)
	}
	reader := multipart.NewReader(bytes.NewReader(body), params["boundary"])
	values := map[string]string{}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			return values
		}
		if err != nil {
			t.Fatalf("NextPart() error = %v", err)
		}
		content, err := io.ReadAll(part)
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}
		values[part.FormName()] = string(content)
	}
}

func TestEnrichTutorMessageTellsTheModelToRejudgeAttachedWork(t *testing.T) {
	// The exact shape a "Revisá" sends: the student's page attached, and the
	// same wording as the previous review, whose "incorrecta" is still in the
	// conversation. Without this instruction the model repeats that verdict
	// even after the student erased and fixed the answer.
	body, contentType := multipartWithImage(t,
		map[string]string{"content": "Revisá mi respuesta actual y ayudame a mejorarla."},
		"image_content", "practice-1.jpg", []byte("fake-png-bytes"),
	)

	rewrittenType, rewrittenBody, err := enrichTutorMessage(contentType, body)
	if err != nil {
		t.Fatalf("enrichTutorMessage() error = %v", err)
	}

	values := readMultipartValues(t, rewrittenType, rewrittenBody)
	if !strings.Contains(values["content"], attachedWorkInstruction) {
		t.Fatalf("content = %q, want the attached-work instruction", values["content"])
	}
	if values["image_content"] != "fake-png-bytes" {
		t.Fatalf("image_content = %q, want the attachment forwarded untouched", values["image_content"])
	}
}

func TestEnrichTutorMessageLeavesTextOnlyMessagesAlone(t *testing.T) {
	body, contentType := multipartBody(t, map[string]string{
		"content": "Dame una pista sin revelar la respuesta.",
	})

	rewrittenType, rewrittenBody, err := enrichTutorMessage(contentType, body)
	if err != nil {
		t.Fatalf("enrichTutorMessage() error = %v", err)
	}

	values := readMultipartValues(t, rewrittenType, rewrittenBody)
	if strings.Contains(values["content"], attachedWorkInstruction) {
		t.Fatal("a message with no attachment should not talk about an attached image")
	}
}

func TestEnrichTutorMessageHandlesImageBeforeContent(t *testing.T) {
	// The instruction depends on a part that may arrive after the one it
	// changes, which is why the parts are buffered before being rewritten.
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	file, err := writer.CreateFormFile("image_content", "page.png")
	if err != nil {
		t.Fatalf("CreateFormFile() error = %v", err)
	}
	if _, err := file.Write([]byte("bytes")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := writer.WriteField("content", "Revisá mi respuesta actual y ayudame a mejorarla."); err != nil {
		t.Fatalf("WriteField() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	rewrittenType, rewrittenBody, err := enrichTutorMessage(writer.FormDataContentType(), body.Bytes())
	if err != nil {
		t.Fatalf("enrichTutorMessage() error = %v", err)
	}

	values := readMultipartValues(t, rewrittenType, rewrittenBody)
	if !strings.Contains(values["content"], attachedWorkInstruction) {
		t.Fatal("the image was not noticed because it arrived after the content part")
	}
}

func multipartWithImage(t *testing.T, fields map[string]string, fileField, filename string, data []byte) ([]byte, string) {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for name, value := range fields {
		if err := writer.WriteField(name, value); err != nil {
			t.Fatalf("WriteField() error = %v", err)
		}
	}
	file, err := writer.CreateFormFile(fileField, filename)
	if err != nil {
		t.Fatalf("CreateFormFile() error = %v", err)
	}
	if _, err := file.Write(data); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	return body.Bytes(), writer.FormDataContentType()
}
