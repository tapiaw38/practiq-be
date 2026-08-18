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
		requesterID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)
		var input struct {
			Title       string `json:"title" binding:"required"`
			Description string `json:"description"`
			Level       int    `json:"level"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "internal server error"})
			return
		}
		out, appErr := uc.Execute(c, requesterID, isSuperAdmin, ucNB.CreateInput{
			CourseID:    courseID,
			Title:       input.Title,
			Description: input.Description,
			Level:       input.Level,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusCreated, out)
	}
}
