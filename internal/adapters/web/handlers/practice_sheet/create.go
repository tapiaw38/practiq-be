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
	// ScheduledAt is RFC 3339; empty clears the schedule.
	ScheduledAt string `json:"scheduled_at"`
	// AvailableUntil is RFC 3339; empty leaves the window open.
	AvailableUntil string `json:"available_until"`
}

func NewCreateHandler(uc ucPS.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		requesterID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		if validationErr := validateSheetTypeAndTestStyle(input.SheetType, input.TestStyle); validationErr != "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": validationErr})
			return
		}

		scheduledAt, err := parseScheduledAt(input.ScheduledAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		availableUntil, err := parseScheduledAt(input.AvailableUntil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		if availableUntil != nil && scheduledAt == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "practice_sheet:invalid-window",
				"message": "the closing date requires an opening date",
			})
			return
		}
		// A window that closes before it opens locks the students out with no
		// sign that anything is wrong.
		if scheduledAt != nil && availableUntil != nil && !availableUntil.After(*scheduledAt) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "practice_sheet:invalid-window",
				"message": "the closing date must be after the opening date",
			})
			return
		}

		output, appErr := uc.Execute(c, requesterID, isSuperAdmin, ucPS.CreateInput{
			CourseID:       courseID,
			SheetType:      input.SheetType,
			TestStyle:      input.TestStyle,
			TopicID:        input.TopicID,
			StrategyID:     input.StrategyID,
			Title:          input.Title,
			Level:          input.Level,
			ExerciseIDs:    input.ExerciseIDs,
			ScheduledAt:    scheduledAt,
			AvailableUntil: availableUntil,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
