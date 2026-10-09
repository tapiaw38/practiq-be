package subject

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucSubject "github.com/tapiaw38/practiq-be/internal/usecases/subject"
)

func NewCreateHandler(uc ucSubject.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucSubject.CreateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, middlewares.GetUserID(c), c.GetHeader("X-School-ID"), middlewares.IsSuperAdmin(c), input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
