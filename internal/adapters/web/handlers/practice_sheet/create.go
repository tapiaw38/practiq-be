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
		requesterID := middlewares.GetUserID(c)
		isAdmin := middlewares.HasRole(c, "admin", "superadmin")
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		if validationErr := validateSheetTypeAndTestStyle(input.SheetType, input.TestStyle); validationErr != "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": validationErr})
			return
		}

		output, appErr := uc.Execute(c, requesterID, isAdmin, ucPS.CreateInput{
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
