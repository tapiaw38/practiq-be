package practicesheet

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucPS "github.com/tapiaw38/practiq-be/internal/usecases/practice_sheet"
)

type createInput struct {
	TopicID     string   `json:"topic_id"`
	StrategyID  string   `json:"strategy_id"`
	Title       string   `json:"title" binding:"required"`
	Level       int      `json:"level"`
	SheetType   string   `json:"sheet_type"`
	TestStyle   string   `json:"test_style"`
	ExerciseIDs []string `json:"exercise_ids"`
}

func NewCreateHandler(uc ucPS.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, ucPS.CreateInput{
			CourseID:    courseID,
			SheetType:   input.SheetType,
			TestStyle:   input.TestStyle,
			TopicID:     input.TopicID,
			StrategyID:  input.StrategyID,
			Title:       input.Title,
			Level:       input.Level,
			ExerciseIDs: input.ExerciseIDs,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}

func NewListHandler(uc ucPS.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		output, appErr := uc.Execute(c, courseID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewGetHandler(uc ucPS.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		output, appErr := uc.Execute(c, id)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

type updateSheetInput struct {
	Title   string `json:"title" binding:"required"`
	TopicID string `json:"topic_id"`
}

func NewUpdateHandler(uc ucPS.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input updateSheetInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, id, ucPS.UpdateInput{
			Title:   input.Title,
			TopicID: input.TopicID,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewDeleteHandler(uc ucPS.DeleteUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if appErr := uc.Execute(c, id); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "practice sheet deleted"})
	}
}

type submitInput struct {
	Attempts []ucPS.AttemptInput `json:"attempts"`
}

func NewSubmitHandler(uc ucPS.SubmitUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		studentID := middlewares.GetUserID(c)
		var input submitInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, id, studentID, ucPS.SubmitInput{
			Attempts: input.Attempts,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
