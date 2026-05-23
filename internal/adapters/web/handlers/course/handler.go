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
		output, appErr := uc.Execute(c, ucCourse.CreateInput{
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

func NewListHandler(uc ucCourse.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := middlewares.GetUserID(c)
		role := c.Query("role")

		input := ucCourse.ListInput{}
		if role == "teacher" {
			input.TeacherID = userID
		} else {
			input.StudentID = userID
		}

		output, appErr := uc.Execute(c, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewGetHandler(uc ucCourse.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		output, appErr := uc.Execute(c, id)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

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
		var input updateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, id, ucCourse.UpdateInput{
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

func NewDeleteHandler(uc ucCourse.DeleteUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if appErr := uc.Execute(c, id); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "course deleted"})
	}
}
