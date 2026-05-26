package studentprogress

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucProgress "github.com/tapiaw38/practiq-be/internal/usecases/student_progress"
)

func NewGetMyProgressHandler(uc ucProgress.GetMyProgressUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := middlewares.GetUserID(c)
		output, appErr := uc.Execute(c, studentID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewGetCourseProgressHandler(uc ucProgress.GetCourseProgressUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := middlewares.GetUserID(c)
		courseID := c.Param("id")
		output, appErr := uc.Execute(c, studentID, courseID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewGetStudentProgressHandler(uc ucProgress.GetStudentProgressUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := c.Param("studentId")
		output, appErr := uc.Execute(c, studentID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func NewGetStudentCourseProgressHandler(uc ucProgress.GetStudentCourseProgressUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := c.Param("studentId")
		courseID := c.Param("courseId")
		output, appErr := uc.Execute(c, studentID, courseID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func NewGetStudentAttemptsHandler(uc ucProgress.GetStudentAttemptsUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := c.Param("studentId")
		sheetID := c.Query("sheet_id")
		if sheetID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "sheet_id query param is required"})
			return
		}
		output, appErr := uc.Execute(c, studentID, sheetID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
