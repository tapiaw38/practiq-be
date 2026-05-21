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

func NewGetMessagesHandler(uc ucAI.GetMessagesUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		conversationID := c.Param("id")
		output, appErr := uc.Execute(c, conversationID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

type helpInput struct {
	ExerciseID string `json:"exercise_id"`
	Question   string `json:"question" binding:"required"`
	HelpType   string `json:"help_type"`
}

func NewHelpHandler(uc ucAI.HelpUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := middlewares.GetUserID(c)
		var input helpInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, ucAI.HelpInput{
			StudentID:  studentID,
			ExerciseID: input.ExerciseID,
			Question:   input.Question,
			HelpType:   input.HelpType,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
