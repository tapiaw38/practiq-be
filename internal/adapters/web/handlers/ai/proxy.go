package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"
)

type requestTransformer func(contentType string, body []byte) (string, []byte, error)

type textMessageInput struct {
	Content string `json:"content"`
	Context string `json:"context"`
}

const maxAssistantProxyRequestBytes = 25 << 20

func proxyToAssistant(uc ucAI.ProxyUsecase, pathBuilder func(*gin.Context) string, transform requestTransformer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxAssistantProxyRequestBytes)
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			if _, ok := err.(*http.MaxBytesError); ok {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{"code": "common:payload-too-large", "message": "assistant attachment exceeds 25 MiB"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "invalid request body"})
			return
		}
		logAssistantProxyBody(c, body)
		contentType := c.GetHeader("Content-Type")
		if transform != nil {
			contentType, body, err = transform(contentType, body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "invalid assistant message"})
				return
			}
		}

		output, appErr := uc.Execute(c, ucAI.ProxyInput{
			UserID:      middlewares.GetUserID(c),
			Method:      c.Request.Method,
			Path:        pathBuilder(c),
			ContentType: contentType,
			Body:        body,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		if output.ContentType != "" {
			c.Header("Content-Type", output.ContentType)
		}
		c.Data(output.StatusCode, output.ContentType, output.Body)
	}
}

const tutorInstruction = `INSTRUCCIONES OBLIGATORIAS DEL ASISTENTE PARA PRACTIQ:
Ayuda al alumno a aprender, no a copiar respuestas. No des respuestas finales ni resuelvas completamente ejercicios evaluables. Da una pista, explicación breve, pregunta guía o siguiente paso.
No cites respuestas correctas, correcciones ni feedback previos que vengan en el contexto o en la imagen.
Si alumno insiste en resultado final, recházalo con amabilidad y guía el procedimiento. Responde en español.

Mensaje del alumno:`

// Added only when the message carries the student's page, because that image
// is the whole answer in a handwritten exercise and the conversation it
// arrives in usually holds an earlier verdict about an earlier version of it.
//
// A student who writes a wrong answer, asks for a review, erases it and fixes
// it was being told "incorrecta" a second time: same wording, no fresh
// context to notice, and its own confident verdict sitting in the history.
// The image is the only thing that changed, so it has to be what decides.
const attachedWorkInstruction = `La imagen adjunta es el estado ACTUAL de la hoja del alumno y reemplaza cualquier imagen o veredicto anterior de esta conversación. El alumno pudo haber borrado y corregido desde tu última respuesta.
Evalúa únicamente lo que ves en esta imagen, partiendo de cero, sin dar por válido ningún juicio previo tuyo. En ejercicios manuscritos la respuesta está solamente en la imagen: es la fuente principal.`

func enrichTutorMessage(contentType string, body []byte) (string, []byte, error) {
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil || !strings.EqualFold(mediaType, "multipart/form-data") {
		return contentType, body, err
	}

	boundary := params["boundary"]
	if boundary == "" {
		return "", nil, io.ErrUnexpectedEOF
	}

	// Buffered rather than streamed because the instruction depends on whether
	// an image is attached, and the parts arrive in whatever order the browser
	// built the form — "content" is usually first, the image after it.
	type formPart struct {
		header textproto.MIMEHeader
		name   string
		file   string
		data   []byte
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
		instruction = tutorInstruction + "\n" + attachedWorkInstruction
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

func logAssistantProxyBody(c *gin.Context, body []byte) {
	contentType := c.GetHeader("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		log.Printf("[assistant_proxy] method=%s path=%s content_type=%q body_bytes=%d", c.Request.Method, c.Request.URL.RequestURI(), contentType, len(body))
		return
	}

	req, err := http.NewRequest(c.Request.Method, c.Request.URL.String(), bytes.NewReader(body))
	if err != nil {
		log.Printf("[assistant_proxy] method=%s path=%s multipart_parse_request_error=%v body_bytes=%d", c.Request.Method, c.Request.URL.RequestURI(), err, len(body))
		return
	}
	req.Header.Set("Content-Type", contentType)
	if err := req.ParseMultipartForm(int64(len(body) + 1024)); err != nil {
		log.Printf("[assistant_proxy] method=%s path=%s multipart_parse_error=%v body_bytes=%d", c.Request.Method, c.Request.URL.RequestURI(), err, len(body))
		return
	}

	imageCount := 0
	imageBytes := int64(0)
	imageNames := []string{}
	audioCount := 0
	audioBytes := int64(0)
	audioNames := []string{}
	if req.MultipartForm != nil {
		for field, files := range req.MultipartForm.File {
			for _, file := range files {
				switch field {
				case "image_content":
					imageCount++
					imageBytes += file.Size
					imageNames = append(imageNames, file.Filename)
				case "voice_content":
					audioCount++
					audioBytes += file.Size
					audioNames = append(audioNames, file.Filename)
				default:
					continue
				}
			}
		}
	}

	log.Printf("[assistant_proxy] method=%s path=%s content_type=%q body_bytes=%d image_count=%d image_bytes=%d image_names=%v audio_count=%d audio_bytes=%d audio_names=%v",
		c.Request.Method,
		c.Request.URL.RequestURI(),
		contentType,
		len(body),
		imageCount,
		imageBytes,
		imageNames,
		audioCount,
		audioBytes,
		audioNames,
	)
}

func NewProxyListConversationsHandler(uc ucAI.ProxyUsecase) gin.HandlerFunc {
	return proxyToAssistant(uc, func(c *gin.Context) string {
		return "/conversation/user"
	}, nil)
}

func NewProxyGetConversationHandler(uc ucAI.ProxyUsecase) gin.HandlerFunc {
	return proxyToAssistant(uc, func(c *gin.Context) string {
		return "/conversation/" + c.Param("id")
	}, nil)
}

func NewProxyCreateConversationHandler(uc ucAI.ProxyUsecase) gin.HandlerFunc {
	return proxyToAssistant(uc, func(c *gin.Context) string {
		return "/conversation/"
	}, nil)
}

func NewProxySendMessageHandler(uc ucAI.ProxyUsecase) gin.HandlerFunc {
	return proxyToAssistant(uc, func(c *gin.Context) string {
		path := "/conversation/" + c.Param("id") + "/message"
		if rawQuery := c.Request.URL.RawQuery; rawQuery != "" {
			path += "?" + rawQuery
		}
		return path
	}, enrichTutorMessage)
}

func NewProxySendTextMessageHandler(uc ucAI.ProxyUsecase) gin.HandlerFunc {
	return proxyToAssistant(uc, func(c *gin.Context) string {
		path := "/conversation/" + c.Param("id") + "/message"
		if rawQuery := c.Request.URL.RawQuery; rawQuery != "" {
			path += "?" + rawQuery
		}
		return path
	}, textMessageToTutorMultipart)
}
