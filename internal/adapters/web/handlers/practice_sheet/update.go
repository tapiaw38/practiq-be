package practicesheet

import (
	"net/http"

	ucPS "github.com/tapiaw38/practiq-be/internal/usecases/practice_sheet"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

type updateSheetInput struct {
	Title       string   `json:"title" binding:"required"`
	TopicID     string   `json:"topic_id"`
	Level       int      `json:"level"`
	SheetType   string   `json:"sheet_type"`
	TestStyle   string   `json:"test_style"`
	ExerciseIDs []string `json:"exercise_ids"`
	// ScheduledAt is RFC 3339; empty clears the schedule.
	ScheduledAt string `json:"scheduled_at"`
	// AvailableUntil is RFC 3339; empty leaves the window open.
	AvailableUntil string `json:"available_until"`
}

func NewUpdateHandler(uc ucPS.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		requesterID := middlewares.GetUserID(c)
		isAdmin := middlewares.HasRole(c, "admin", "superadmin")

		var input updateSheetInput
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

		output, appErr := uc.Execute(c, requesterID, isAdmin, id, ucPS.UpdateInput{
			Title:          input.Title,
			TopicID:        input.TopicID,
			Level:          input.Level,
			SheetType:      input.SheetType,
			TestStyle:      input.TestStyle,
			ExerciseIDs:    input.ExerciseIDs,
			ScheduledAt:    scheduledAt,
			AvailableUntil: availableUntil,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
