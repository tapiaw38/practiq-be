package subject

import (
	"net/http"

	ucSubject "github.com/tapiaw38/practiq-be/internal/usecases/subject"

	"github.com/gin-gonic/gin"
)

func NewListHandler(uc ucSubject.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
