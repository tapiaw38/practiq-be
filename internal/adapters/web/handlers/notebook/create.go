package notebook

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"
)

func NewCreateHandler(uc ucNB.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		teacherID := middlewares.GetUserID(c)
		var input struct {
			Title       string `json:"title" binding:"required"`
			Description string `json:"description"`
			Level       int    `json:"level"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		out, err := uc.Execute(c, teacherID, ucNB.CreateInput{
			CourseID:    courseID,
			Title:       input.Title,
			Description: input.Description,
			Level:       input.Level,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, out)
	}
}
