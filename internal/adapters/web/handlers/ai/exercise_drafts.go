package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"
	ucExercise "github.com/tapiaw38/practiq-be/internal/usecases/exercise"
)

// Drafts are deliberately not persisted. A teacher reviews every field before
// the normal exercise creation endpoint writes anything to the course.
const maxExerciseSourceBytes = 20 << 20

func NewExerciseDraftsHandler(uc ucAI.ProxyUsecase, exercises ucExercise.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// List uses the same requesterCanWriteTopic guard as create. Do this
		// before accepting bytes, so a teacher cannot spend AI quota on another
		// school's topic.
		if _, appErr := exercises.Execute(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), ucExercise.ListInput{TopicID: c.Param("id"), Limit: 1}); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxExerciseSourceBytes)
		if err := c.Request.ParseMultipartForm(maxExerciseSourceBytes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "Seleccioná un PDF, DOCX o imagen de hasta 20 MB."})
			return
		}
		file, header, err := c.Request.FormFile("source")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "Falta el archivo fuente."})
			return
		}
		defer file.Close()
		content, err := io.ReadAll(io.LimitReader(file, maxExerciseSourceBytes+1))
		if err != nil || len(content) > maxExerciseSourceBytes {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"code": "common:payload-too-large", "message": "El archivo supera 20 MB."})
			return
		}
		ext := strings.ToLower(filepath.Ext(header.Filename))
		isDocument := ext == ".pdf" || ext == ".docx"
		isImage := ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".webp"
		if !isDocument && !isImage {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"code": "common:unsupported-media", "message": "Formatos permitidos: PDF, DOCX, PNG, JPG o WebP."})
			return
		}

		userID := middlewares.GetUserID(c)
		conversationID, appErr := createDraftConversation(c, uc, userID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		body, contentType, err := draftMessageBody(content, header.Filename, isDocument, c.PostForm("count"), c.PostForm("difficulty"), c.PostForm("instruction"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "No se pudo preparar el documento."})
			return
		}
		response, appErr := uc.Execute(c, ucAI.ProxyInput{UserID: userID, Method: http.MethodPost, Path: "/conversation/" + conversationID + "/message", ContentType: contentType, Body: body})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			c.Data(response.StatusCode, response.ContentType, response.Body)
			return
		}
		drafts, err := parseDrafts(response.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"code": "ai:invalid-drafts", "message": "La IA no devolvió ejercicios utilizables. Intentá nuevamente."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": drafts})
	}
}

func createDraftConversation(c *gin.Context, uc ucAI.ProxyUsecase, userID string) (string, apperrors.ApplicationError) {
	// Kept in Gillie only to use its existing document/OCR pipeline. Its title
	// makes these ephemeral authoring conversations identifiable for cleanup.
	payload := []byte(`{"title":"Practiq · borrador de ejercicios","is_sandbox":true}`)
	response, appErr := uc.Execute(c, ucAI.ProxyInput{UserID: userID, Method: http.MethodPost, Path: "/conversation/", ContentType: "application/json", Body: payload})
	if appErr != nil {
		return "", appErr
	}
	var out struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 || json.Unmarshal(response.Body, &out) != nil || out.Data.ID == "" {
		return "", apperrors.NewInternalError(fmt.Errorf("assistant conversation could not be created"))
	}
	return out.Data.ID, nil
}

func draftMessageBody(content []byte, filename string, document bool, count, difficulty, instruction string) ([]byte, string, error) {
	if count == "" {
		count = "5"
	}
	if difficulty == "" {
		difficulty = "5"
	}
	prompt := fmt.Sprintf(`Creá exactamente %s ejercicios en español argentino, dificultad %s, basados solamente en archivo adjunto. Devolvé ÚNICAMENTE JSON válido, sin markdown: {"drafts":[{"type":"open_text|multiple_choice|equation","question":"...","correct_answer":"...","explanation":"...","difficulty":1,"metadata":{"options":["..."]}}]}. Usá multiple_choice solo con cuatro opciones y respuesta correcta incluida. No generes canvas, handwritten, attachment ni fill_blanks. %s`, count, difficulty, strings.TrimSpace(instruction))
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	if err := w.WriteField("content", prompt); err != nil {
		return nil, "", err
	}
	field := "image_content"
	if document {
		field = "document_content"
	}
	p, err := w.CreateFormFile(field, filename)
	if err != nil {
		return nil, "", err
	}
	if _, err = p.Write(content); err != nil {
		return nil, "", err
	}
	if err = w.Close(); err != nil {
		return nil, "", err
	}
	return b.Bytes(), w.FormDataContentType(), nil
}

func parseDrafts(body []byte) ([]map[string]interface{}, error) {
	var envelope struct {
		Data []struct {
			Content string `json:"content"`
			Sender  string `json:"sender"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	for i := len(envelope.Data) - 1; i >= 0; i-- {
		if envelope.Data[i].Sender == "assistant" {
			text := strings.TrimSpace(envelope.Data[i].Content)
			text = strings.TrimPrefix(strings.TrimSuffix(text, "```"), "```json")
			var output struct {
				Drafts []map[string]interface{} `json:"drafts"`
			}
			if err := json.Unmarshal([]byte(strings.TrimSpace(text)), &output); err != nil {
				return nil, err
			}
			if len(output.Drafts) == 0 {
				return nil, fmt.Errorf("empty drafts")
			}
			return output.Drafts, nil
		}
	}
	return nil, fmt.Errorf("assistant response missing")
}
