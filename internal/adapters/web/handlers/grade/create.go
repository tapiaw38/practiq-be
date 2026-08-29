package grade

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucGrade "github.com/tapiaw38/practiq-be/internal/usecases/grade"
)

type createInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	VisualTheme string `json:"visual_theme"`
}

func NewCreateHandler(uc ucGrade.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, ucGrade.CreateInput{
			Name:        input.Name,
			Description: input.Description,
			VisualTheme: input.VisualTheme,
			CreatedBy:   middlewares.GetUserID(c),
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
