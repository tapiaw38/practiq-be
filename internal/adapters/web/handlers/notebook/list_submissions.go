package notebook

import (
	"net/http"
	"strings"

	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewListSubmissionsHandler(uc ucNB.ListSubmissionsUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		reviewed := ""
		if raw, ok := c.GetQuery("reviewed"); ok {
			if raw == "true" {
				reviewed = "reviewed"
			} else if raw == "false" {
				reviewed = "unreviewed"
			}
		}
		teacherID := ""
		if !middlewares.HasRole(c, "admin", "superadmin") {
			teacherID = middlewares.GetUserID(c)
		}
		output, err := uc.Execute(c, ucNB.ListSubmissionsInput{
			NotebookID: strings.TrimSpace(c.Query("notebook_id")),
			StudentID:  strings.TrimSpace(c.Query("student_id")),
			CourseID:   strings.TrimSpace(c.Query("course_id")),
			Reviewed:   reviewed,
			TeacherID:  teacherID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "notebook:list-submissions-error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
