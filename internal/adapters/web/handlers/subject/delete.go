package subject

import (
	"net/http"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucSubject "github.com/tapiaw38/practiq-be/internal/usecases/subject"

	"github.com/gin-gonic/gin"
)

func NewDeleteHandler(uc ucSubject.DeleteUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if appErr := uc.Execute(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), id); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "subject deleted"})
	}
}
