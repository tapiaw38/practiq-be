package exercise

import (
	"net/http"

	ucExercise "github.com/tapiaw38/practiq-be/internal/usecases/exercise"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

type updateInput struct {
	Type          string `json:"type"`
	Question      string `json:"question"`
	CorrectAnswer string `json:"correct_answer"`
	Explanation   string `json:"explanation"`
	Difficulty    int    `json:"difficulty"`
	Metadata      string `json:"metadata"`
}

func NewUpdateHandler(uc ucExercise.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		requesterID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)

		var input updateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, requesterID, isSuperAdmin, id, ucExercise.UpdateInput{
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

		c.JSON(http.StatusOK, output)
	}
}
