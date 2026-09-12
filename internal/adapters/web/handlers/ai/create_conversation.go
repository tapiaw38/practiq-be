package ai

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"
)

type createConversationInput struct {
	CourseID        string `json:"course_id"`
	PracticeSheetID string `json:"practice_sheet_id"`
}

func NewCreateConversationHandler(uc ucAI.CreateConversationUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := middlewares.GetUserID(c)
		var input createConversationInput
		c.ShouldBindJSON(&input)

		output, appErr := uc.Execute(c, ucAI.CreateConversationInput{
			StudentID:       studentID,
			CourseID:        input.CourseID,
			PracticeSheetID: input.PracticeSheetID,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
