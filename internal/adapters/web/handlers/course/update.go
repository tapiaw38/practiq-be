package course

import (
	"net/http"

	ucCourse "github.com/tapiaw38/practiq-be/internal/usecases/course"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

type updateInput struct {
	GradeID     string `json:"grade_id"`
	SubjectID   string `json:"subject_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	Subject     string `json:"subject"`
}

func NewUpdateHandler(uc ucCourse.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		requesterID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)

		var input updateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, requesterID, isSuperAdmin, id, ucCourse.UpdateInput{
			GradeID:     input.GradeID,
			SubjectID:   input.SubjectID,
			Title:       input.Title,
			Description: input.Description,
			Level:       input.Level,
			Subject:     input.Subject,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
