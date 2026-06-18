package course

import (
	"net/http"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucCourse "github.com/tapiaw38/practiq-be/internal/usecases/course"

	"github.com/gin-gonic/gin"
)

func NewGetHandler(uc ucCourse.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID := middlewares.GetUserID(c)
		isAdmin := middlewares.HasRole(c, "admin", "superadmin")
		output, appErr := uc.Execute(c, userID, isAdmin, id)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
