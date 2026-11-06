package studentreport

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-be/internal/domain"
	ucReport "github.com/tapiaw38/practiq-be/internal/usecases/student_report"
)

func NewGeneratePDFHandler(uc ucReport.GeneratePDFUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		teacherID := middlewares.GetUserID(c)
		studentID := c.Param("studentId")
		isSuperAdmin := middlewares.IsSuperAdmin(c)

		if studentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "studentId is required"})
			return
		}

		filter := domain.StudentReportFilter{
			StudentID: studentID,
			CourseID:  c.Query("course_id"),
		}

		if fromStr := c.Query("from"); fromStr != "" {
			t, err := time.Parse("2006-01-02", fromStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "invalid from date format, use YYYY-MM-DD"})
				return
			}
			filter.From = &t
		}

		if toStr := c.Query("to"); toStr != "" {
			t, err := time.Parse("2006-01-02", toStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "invalid to date format, use YYYY-MM-DD"})
				return
			}

			t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filter.To = &t
		}

		output, appErr := uc.Execute(c, teacherID, isSuperAdmin, filter, c.GetHeader("Authorization"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		studentPrefix := studentID
		if len(studentID) > 8 {
			studentPrefix = studentID[:8]
		}
		filename := fmt.Sprintf("reporte_progreso_%s_%s.pdf", studentPrefix, time.Now().Format("20060102"))

		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Header("Content-Length", fmt.Sprintf("%d", len(output.Content)))

		c.Data(http.StatusOK, "application/pdf", output.Content)
	}
}
