package notebook

import (
	"log"
	"net/http"
	"strconv"
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

		input := ucNB.ListSubmissionsInput{
			NotebookID: strings.TrimSpace(c.Query("notebook_id")),
			StudentID:  strings.TrimSpace(c.Query("student_id")),
			CourseID:   strings.TrimSpace(c.Query("course_id")),
			Reviewed:   reviewed,
			TeacherID:  teacherID,
		}

		if limitStr := c.Query("limit"); limitStr != "" {
			if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
				input.Limit = limit
			}
		}

		if offsetStr := c.Query("offset"); offsetStr != "" {
			if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
				input.Offset = offset
			}
		}

		output, err := uc.Execute(c, input)
		if err != nil {
			log.Printf("[notebook handler] list submissions error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "notebook:list-submissions-error", "message": "internal server error"})
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
