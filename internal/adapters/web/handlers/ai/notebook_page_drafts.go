package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"
	ucNotebook "github.com/tapiaw38/practiq-be/internal/usecases/notebook"
)

// NewNotebookPageDraftsHandler creates text-page drafts only. The teacher
// reviews them and the normal page endpoint is still the only persistent write.
func NewNotebookPageDraftsHandler(proxy ucAI.ProxyUsecase, notebooks ucNotebook.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := middlewares.GetUserID(c)
		notebook, appErr := notebooks.Execute(c, userID, middlewares.IsSuperAdmin(c), c.Param("id"), "")
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		if notebook == nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "notebook:not-found", "message": "Cuaderno no encontrado."})
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxExerciseSourceBytes)
		if err := c.Request.ParseMultipartForm(maxExerciseSourceBytes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "Seleccioná un PDF, DOCX o imagen de hasta 20 MB."})
			return
		}
		instruction := strings.TrimSpace(c.PostForm("instruction"))
		content, filename, document, err := notebookPageSource(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		if len(content) == 0 && instruction == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "Subí un archivo o escribí qué hoja querés crear."})
			return
		}

		conversationID, appErr := createDraftConversation(c, proxy, userID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		body, contentType, err := notebookPageDraftMessageBody(content, filename, document, c.PostForm("count"), instruction)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "No se pudo preparar el pedido."})
			return
		}
		response, appErr := proxy.Execute(c, ucAI.ProxyInput{UserID: userID, Method: http.MethodPost, Path: "/conversation/" + conversationID + "/message", ContentType: contentType, Body: body})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			c.Data(response.StatusCode, response.ContentType, response.Body)
			return
		}
		pages, err := parseNotebookPageDrafts(response.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"code": "ai:invalid-page-drafts", "message": "La IA no devolvió hojas utilizables. Intentá nuevamente."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": pages})
	}
}

func notebookPageSource(c *gin.Context) ([]byte, string, bool, error) {
	file, header, err := c.Request.FormFile("source")
	if err != nil {
		return nil, "", false, nil
	}
	defer file.Close()
	content, err := io.ReadAll(io.LimitReader(file, maxExerciseSourceBytes+1))
	if err != nil || len(content) > maxExerciseSourceBytes {
		return nil, "", false, fmt.Errorf("el archivo supera 20 MB o no se pudo leer")
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	document := ext == ".pdf" || ext == ".docx"
	image := ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".webp"
	if !document && !image {
		return nil, "", false, fmt.Errorf("formatos permitidos: PDF, DOCX, PNG, JPG o WebP")
	}
	return content, header.Filename, document, nil
}

func notebookPageDraftMessageBody(content []byte, filename string, document bool, rawCount, instruction string) ([]byte, string, error) {
	count, err := strconv.Atoi(rawCount)
	if err != nil || count < 1 || count > 5 {
		count = 1
	}
	source := "sobre la indicación del docente"
	if len(content) > 0 {
		source = "basadas solamente en el archivo adjunto"
	}
	prompt := fmt.Sprintf(`Creá exactamente %d hojas de cuaderno educativas en español argentino, %s. Devolvé ÚNICAMENTE JSON válido, sin markdown: {"pages":[{"title":"...","content_type":"text","content_data":"...","instructions":"..."}]}. Cada hoja debe ser autocontenida, clara para estudiante, con explicación breve, ejemplo cuando ayude y una actividad final concreta. content_type debe ser siempre "text". content_data debe estar listo para mostrar, con párrafos y saltos de línea; no uses Markdown, HTML ni respuestas de ejercicios. instructions debe indicar qué debe hacer el alumno. %s`, count, source, instruction)
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	if err := w.WriteField("content", prompt); err != nil {
		return nil, "", err
	}
	if len(content) > 0 {
		field := "image_content"
		if document {
			field = "document_content"
		}
		part, err := w.CreateFormFile(field, filename)
		if err != nil {
			return nil, "", err
		}
		if _, err = part.Write(content); err != nil {
			return nil, "", err
		}
	}
	if err := w.Close(); err != nil {
		return nil, "", err
	}
	return b.Bytes(), w.FormDataContentType(), nil
}

func parseNotebookPageDrafts(body []byte) ([]map[string]string, error) {
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
		if envelope.Data[i].Sender != "assistant" {
			continue
		}
		text := strings.TrimSpace(envelope.Data[i].Content)
		text = strings.TrimPrefix(strings.TrimSuffix(text, "```"), "```json")
		var output struct {
			Pages []map[string]string `json:"pages"`
		}
		if err := json.Unmarshal([]byte(strings.TrimSpace(text)), &output); err != nil {
			return nil, err
		}
		if len(output.Pages) == 0 {
			return nil, fmt.Errorf("empty pages")
		}
		for _, page := range output.Pages {
			if strings.TrimSpace(page["title"]) == "" || strings.TrimSpace(page["content_data"]) == "" {
				return nil, fmt.Errorf("page missing title or content")
			}
			page["content_type"] = "text"
		}
		return output.Pages, nil
	}
	return nil, fmt.Errorf("assistant response missing")
}
