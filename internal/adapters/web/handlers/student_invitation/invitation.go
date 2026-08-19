package studentinvitation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucInvitation "github.com/tapiaw38/practiq-be/internal/usecases/student_invitation"
)

type redeemInput struct {
	Code string `json:"code" binding:"required"`
}

// NewCreateHandler genera el código del docente que hace el pedido. El
// teacher_id sale del token: aceptarlo del cuerpo dejaría que un docente
// genere códigos a nombre de otro.
func NewCreateHandler(uc ucInvitation.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, middlewares.GetUserID(c))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}

func NewGetActiveHandler(uc ucInvitation.GetActiveUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, middlewares.GetUserID(c))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewRevokeHandler(uc ucInvitation.RevokeUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if appErr := uc.Execute(c, c.Param("id"), middlewares.GetUserID(c)); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func NewRedeemHandler(uc ucInvitation.RedeemUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input redeemInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, middlewares.GetUserID(c), input.Code)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
