package school

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucSchool "github.com/tapiaw38/practiq-be/internal/usecases/school"
)

func NewAddMemberHandler(uc ucSchool.AddMemberUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucSchool.AddMemberInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		if appErr := uc.Execute(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id"), input); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}
