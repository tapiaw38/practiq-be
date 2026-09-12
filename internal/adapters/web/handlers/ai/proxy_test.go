package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
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
