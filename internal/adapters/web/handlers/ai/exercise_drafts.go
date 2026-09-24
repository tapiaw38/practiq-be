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
		// A source file is one way to say what the exercises should be about;
		// a written instruction is the other. Requiring the file meant a
		// teacher who could describe the topic in a sentence had to find a
		// document first. One of the two has to be there: with neither, there
		// is nothing to base the exercises on.
		instruction := strings.TrimSpace(c.PostForm("instruction"))
		var (
			content    []byte
			filename   string
			isDocument bool
		)
		if file, header, err := c.Request.FormFile("source"); err == nil {
			defer file.Close()
			content, err = io.ReadAll(io.LimitReader(file, maxExerciseSourceBytes+1))
			if err != nil || len(content) > maxExerciseSourceBytes {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{"code": "common:payload-too-large", "message": "El archivo supera 20 MB."})
				return
			}
			ext := strings.ToLower(filepath.Ext(header.Filename))
			isDocument = ext == ".pdf" || ext == ".docx"
			isImage := ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".webp"
			if !isDocument && !isImage {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{"code": "common:unsupported-media", "message": "Formatos permitidos: PDF, DOCX, PNG, JPG o WebP."})
				return
			}
			filename = header.Filename
		} else if instruction == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "common:bad-request",
				"message": "Subí un archivo o escribí sobre qué tema generar los ejercicios.",
			})
			return
		}

		userID := middlewares.GetUserID(c)
		conversationID, appErr := createDraftConversation(c, uc, userID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		body, contentType, err := draftMessageBody(content, filename, isDocument, c.PostForm("count"), c.PostForm("difficulty"), instruction, c.PostForm("exercise_type"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "No se pudo preparar el pedido."})
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
		drafts, err := parseDrafts(response.Body, c.PostForm("exercise_type"))
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

// draftTypeRule turns the teacher's choice into an order the assistant can
// follow. Offered as one option among four in a union, fill_blanks never came
// back: the model answered with the simplest shape that satisfied the prompt.
func draftTypeRule(exerciseType string) string {
	switch exerciseType {
	case "open_text":
		return ` Generá únicamente ejercicios de tipo "open_text".`
	case "equation":
		return ` Generá únicamente ejercicios de tipo "equation".`
	case "multiple_choice":
		return ` Generá únicamente ejercicios de tipo "multiple_choice".`
	case "fill_blanks":
		return ` Generá únicamente ejercicios de tipo "fill_blanks". Ejemplo del formato exacto: {"type":"fill_blanks","question":"El agua hierve a {{1}} grados y se congela a {{2}} grados.","correct_answer":"","explanation":"...","difficulty":1,"metadata":{"blanks":[{"id":1,"answer":"100"},{"id":2,"answer":"0"}],"distractors":["50","212"],"layout":"text"}}`
	case "canvas":
		return ` Generá únicamente ejercicios de tipo "canvas". La consigna debe pedir al alumno resolver o representar el trabajo dibujando en el lienzo; agregá correct_answer como referencia breve para corregir.`
	case "attachment":
		return ` Generá únicamente ejercicios de tipo "attachment". La consigna debe pedir una entrega concreta y metadata.accept debe ser una lista de formatos entre "audio", "pdf", "image" y "doc"; correct_answer puede quedar vacío.`
	default:
		return " Variá los tipos entre los disponibles."
	}
}

func draftMessageBody(content []byte, filename string, document bool, count, difficulty, instruction, exerciseType string) ([]byte, string, error) {
	if count == "" {
		count = "1"
	}
	if difficulty == "" {
		difficulty = "1"
	}
	source := "basados solamente en el archivo adjunto"
	if len(content) == 0 {
		source = "sobre el tema indicado abajo"
	}
	prompt := fmt.Sprintf(`Creá exactamente %s ejercicios en español argentino, dificultad %s, %s. Devolvé ÚNICAMENTE JSON válido, sin markdown: {"drafts":[{"type":"open_text|multiple_choice|equation|canvas|attachment|fill_blanks","question":"...","correct_answer":"...","explanation":"...","difficulty":1,"metadata":{"options":["..."],"blanks":[{"id":1,"answer":"..."}],"distractors":["..."],"layout":"text","accept":["image"]}}]}. Tipos permitidos: open_text, multiple_choice, equation, canvas, attachment y fill_blanks. Nunca generes handwritten/manuscrito. Usá multiple_choice solo con cuatro opciones en metadata.options y correct_answer igual a una de ellas. Para fill_blanks escribí el enunciado con marcadores {{1}}, {{2}}, cada número una sola vez y en orden; poné una entrada en metadata.blanks por cada marcador con su respuesta, agregá opciones incorrectas en metadata.distractors, usá metadata.layout "code" solo si el enunciado es código, y dejá correct_answer vacío porque se arma con los huecos. Para canvas pedí que el alumno resuelva o represente en el lienzo y escribí correct_answer como referencia de corrección. Para attachment pedí una entrega concreta, usá metadata.accept con uno o más de audio, pdf, image o doc, y podés dejar correct_answer vacío. Para los demás tipos no incluyas blanks, distractors ni accept.%s %s`, count, difficulty, source, draftTypeRule(exerciseType), strings.TrimSpace(instruction))
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
		p, err := w.CreateFormFile(field, filename)
		if err != nil {
			return nil, "", err
		}
		if _, err = p.Write(content); err != nil {
			return nil, "", err
		}
	}
	if err := w.Close(); err != nil {
		return nil, "", err
	}
	return b.Bytes(), w.FormDataContentType(), nil
}

func parseDrafts(body []byte, requestedType string) ([]map[string]interface{}, error) {
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
			for _, draft := range output.Drafts {
				if err := validateDraftShape(draft, requestedType); err != nil {
					return nil, err
				}
			}
			return output.Drafts, nil
		}
	}
	return nil, fmt.Errorf("assistant response missing")
}

// validateDraftShape is a boundary guard, not a replacement for the editor's
// review. It makes the assistant contract match the manual exercise catalogue
// and prevents unsupported or handwritten drafts from reaching the UI.
func validateDraftShape(draft map[string]interface{}, requestedType string) error {
	typ, _ := draft["type"].(string)
	question, _ := draft["question"].(string)
	allowed := map[string]bool{
		"open_text": true, "multiple_choice": true, "equation": true,
		"canvas": true, "attachment": true, "fill_blanks": true,
	}
	if !allowed[typ] {
		return fmt.Errorf("unsupported draft type %q", typ)
	}
	if requestedType != "" && typ != requestedType {
		return fmt.Errorf("draft type %q does not match requested type %q", typ, requestedType)
	}
	if strings.TrimSpace(question) == "" {
		return fmt.Errorf("draft question is empty")
	}
	return nil
}
