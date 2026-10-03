package learningstrategy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ucLS "github.com/tapiaw38/practiq-be/internal/usecases/learning_strategy"
)

func NewCreateHandler(uc ucLS.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucLS.CreateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
