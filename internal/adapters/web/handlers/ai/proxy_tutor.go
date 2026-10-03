package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/textproto"
	"strings"
)

type formPart struct {
	header textproto.MIMEHeader
	name   string
	file   string
	data   []byte
}

const tutorInstruction = `INSTRUCCIONES OBLIGATORIAS DEL ASISTENTE PARA PRACTIQ:
Ayuda al alumno a aprender, no a copiar respuestas. No des respuestas finales ni resuelvas completamente ejercicios evaluables. Da una pista, explicación breve, pregunta guía o siguiente paso.
No cites respuestas correctas, correcciones ni feedback previos que vengan en el contexto o en la imagen.
Si alumno insiste en resultado final, recházalo con amabilidad y guía el procedimiento. Responde en español.

Mensaje del alumno:`

const attachedWorkInstruction = `La imagen adjunta es el estado ACTUAL de la hoja del alumno y reemplaza cualquier imagen o veredicto anterior de esta conversación. El alumno pudo haber borrado y corregido desde tu última respuesta.
Evalúa únicamente lo que ves en esta imagen, partiendo de cero, sin dar por válido ningún juicio previo tuyo. En ejercicios manuscritos la respuesta está solamente en la imagen: es la fuente principal.`

type textMessageInput struct {
	Content string `json:"content"`
	Context string `json:"context"`
}

func enrichTutorMessage(contentType string, body []byte) (string, []byte, error) {
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil || !strings.EqualFold(mediaType, "multipart/form-data") {
		return contentType, body, err
	}
	boundary := params["boundary"]
	if boundary == "" {
		return "", nil, io.ErrUnexpectedEOF
	}

	reader := multipart.NewReader(bytes.NewReader(body), boundary)
	parts := make([]formPart, 0, 8)
	hasImage := false
	for {
		part, partErr := reader.NextPart()
		if partErr == io.EOF {
			break
		}
		if partErr != nil {
			return "", nil, partErr
		}
		data, readErr := io.ReadAll(part)
		if readErr != nil {
			return "", nil, readErr
		}
		if part.FormName() == "image_content" && len(data) > 0 {
			hasImage = true
		}
		parts = append(parts, formPart{header: part.Header, name: part.FormName(), file: part.FileName(), data: data})
	}

	instruction := tutorInstruction
	if hasImage {
		instruction += "\n" + attachedWorkInstruction
	}
	var rewritten bytes.Buffer
	writer := multipart.NewWriter(&rewritten)
	for _, part := range parts {
		destination, createErr := writer.CreatePart(part.header)
		if createErr != nil {
			return "", nil, createErr
		}
		payload := part.data
		if part.name == "content" && part.file == "" {
			message := strings.TrimSpace(string(part.data))
			if !strings.Contains(message, "POLITICA OBLIGATORIA:") && !strings.Contains(message, "INSTRUCCIONES OBLIGATORIAS DEL ASISTENTE PARA PRACTIQ:") {
				message = instruction + "\n" + message
			}
			payload = []byte(message)
		}
		if _, writeErr := destination.Write(payload); writeErr != nil {
			return "", nil, writeErr
		}
	}
	if err := writer.Close(); err != nil {
		return "", nil, err
	}
	return writer.FormDataContentType(), rewritten.Bytes(), nil
}

func textMessageToTutorMultipart(contentType string, body []byte) (string, []byte, error) {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil || !strings.EqualFold(mediaType, "application/json") {
		return contentType, body, err
	}
	var input textMessageInput
	if err := json.Unmarshal(body, &input); err != nil {
		return "", nil, err
	}
	if strings.TrimSpace(input.Content) == "" {
		return "", nil, io.ErrUnexpectedEOF
	}
	var multipartBody bytes.Buffer
	writer := multipart.NewWriter(&multipartBody)
	if err := writer.WriteField("content", input.Content); err != nil {
		return "", nil, err
	}
	if strings.TrimSpace(input.Context) != "" {
		if err := writer.WriteField("context", input.Context); err != nil {
			return "", nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return "", nil, err
	}
	return enrichTutorMessage(writer.FormDataContentType(), multipartBody.Bytes())
}
