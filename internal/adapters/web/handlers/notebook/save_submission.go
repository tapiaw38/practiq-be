package notebook

import (
	"log"
	"net/http"
	"strings"

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
		if err := uc.Execute(c, ucNB.SaveSubmissionInput{
			PageID:     pageID,
			StudentID:  studentID,
			CanvasData: input.CanvasData,
			AnswerText: input.AnswerText,
		}); err != nil {
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
