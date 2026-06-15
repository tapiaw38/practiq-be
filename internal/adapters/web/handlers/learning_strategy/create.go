package learningstrategy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ucLS "github.com/tapiaw38/practiq-be/internal/usecases/learning_strategy"
)

type createInput struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
}

func NewCreateHandler(uc ucLS.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, ucLS.CreateInput{
			Name:        input.Name,
			Code:        input.Code,
			Description: input.Description,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
