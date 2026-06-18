package exercise

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucExercise "github.com/tapiaw38/practiq-be/internal/usecases/exercise"
)

type createInput struct {
	Type          string `json:"type" binding:"required"`
	Question      string `json:"question" binding:"required"`
	CorrectAnswer string `json:"correct_answer"`
	Explanation   string `json:"explanation"`
	Difficulty    int    `json:"difficulty"`
	Metadata      string `json:"metadata"`
}

func NewCreateHandler(uc ucExercise.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		topicID := c.Param("id")
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		userID := middlewares.GetUserID(c)
		isAdmin := middlewares.HasRole(c, "admin", "superadmin")
		output, appErr := uc.Execute(c, userID, isAdmin, ucExercise.CreateInput{
			TopicID:       topicID,
			Type:          input.Type,
			Question:      input.Question,
			CorrectAnswer: input.CorrectAnswer,
			Explanation:   input.Explanation,
			Difficulty:    input.Difficulty,
			Metadata:      input.Metadata,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
