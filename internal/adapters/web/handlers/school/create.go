package school

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ucSchool "github.com/tapiaw38/practiq-be/internal/usecases/school"
)

func NewCreateHandler(uc ucSchool.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucSchool.CreateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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
