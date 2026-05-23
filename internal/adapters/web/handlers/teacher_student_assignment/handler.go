package teacherstudentassignment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucAssignment "github.com/tapiaw38/practiq-be/internal/usecases/teacher_student_assignment"
)

type assignInput struct {
	TeacherID string `json:"teacher_id" binding:"required"`
	StudentID string `json:"student_id" binding:"required"`
}

func NewAssignHandler(uc ucAssignment.AssignUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input assignInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, input.TeacherID, input.StudentID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func NewUnassignHandler(uc ucAssignment.UnassignUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		teacherID := c.Param("teacherId")
		studentID := c.Param("studentId")
		output, appErr := uc.Execute(c, teacherID, studentID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func NewListStudentsHandler(uc ucAssignment.ListStudentsUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		teacherID := c.Param("teacherId")
		output, appErr := uc.Execute(c, teacherID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func NewListTeachersHandler(uc ucAssignment.ListTeachersUsecase) gin.HandlerFunc {
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

func NewListMyStudentsHandler(uc ucAssignment.ListStudentsUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		teacherID := middlewares.GetUserID(c)
		output, appErr := uc.Execute(c, teacherID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
