package notebook

import (
	"net/http"

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
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		if err := uc.Execute(c, ucNB.SaveSubmissionInput{
			PageID:     pageID,
			StudentID:  studentID,
			CanvasData: input.CanvasData,
			AnswerText: input.AnswerText,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
