package exercise

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

		output, appErr := uc.Execute(c, ucExercise.CreateInput{
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

func NewListHandler(uc ucExercise.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		topicID := c.Param("id")
		output, appErr := uc.Execute(c, topicID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

type updateInput struct {
	Type          string `json:"type"`
	Question      string `json:"question"`
	CorrectAnswer string `json:"correct_answer"`
	Explanation   string `json:"explanation"`
	Difficulty    int    `json:"difficulty"`
}

func NewUpdateHandler(uc ucExercise.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input updateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, id, ucExercise.UpdateInput{
			Type:          input.Type,
			Question:      input.Question,
			CorrectAnswer: input.CorrectAnswer,
			Explanation:   input.Explanation,
			Difficulty:    input.Difficulty,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewDeleteHandler(uc ucExercise.DeleteUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if appErr := uc.Execute(c, id); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "exercise deleted"})
	}
}
