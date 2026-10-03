package school

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucSchool "github.com/tapiaw38/practiq-be/internal/usecases/school"
)

func NewRemoveMemberHandler(uc ucSchool.RemoveMemberUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if appErr := uc.Execute(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id"), c.Param("userId")); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}
