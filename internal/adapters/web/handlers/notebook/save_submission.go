package notebook

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	notebookRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/notebook"
	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewSaveSubmissionHandler(uc ucNB.SaveSubmissionUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageID := c.Param("id")
		studentID := middlewares.GetUserID(c)
		var input struct {
			CanvasData string `json:"canvas_data"`
			AnswerText string `json:"answer_text"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "internal server error"})
			return
		}
		// Versioned like the async path, and on the same clock: the two write
		// to the same row, so a delivery sent here must be ordered against one
		// still being processed there.
		if err := uc.Execute(c, ucNB.SaveSubmissionInput{
			PageID:     pageID,
			StudentID:  studentID,
			CanvasData: input.CanvasData,
			AnswerText: input.AnswerText,
			Version:    time.Now().UTC().UnixNano(),
		}); err != nil {
			// The student already replaced this answer. Nothing to save and
			// nothing to report: the newer one is stored.
			if errors.Is(err, notebookRepo.ErrStaleSubmission) {
				c.JSON(http.StatusNoContent, nil)
				return
			}
			if strings.Contains(err.Error(), "forbidden") {
				c.JSON(http.StatusForbidden, gin.H{"code": "common:forbidden", "message": "forbidden"})
				return
			}
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{"code": "notebook:not-found", "message": err.Error()})
				return
			}
			log.Printf("[notebook handler] save submission error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
