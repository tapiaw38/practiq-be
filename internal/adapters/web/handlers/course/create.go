package course

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucCourse "github.com/tapiaw38/practiq-be/internal/usecases/course"
)

type createInput struct {
	GradeID     string `json:"grade_id" binding:"required"`
	SubjectID   string `json:"subject_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Level       string `json:"level"`
	Subject     string `json:"subject"`
}

func NewCreateHandler(uc ucCourse.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		userID := middlewares.GetUserID(c)
		output, appErr := uc.Execute(c, middlewares.HasRole(c, "teacher", "admin", "superadmin"), ucCourse.CreateInput{
			TeacherID:   userID,
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

		c.JSON(http.StatusCreated, output)
	}
}
