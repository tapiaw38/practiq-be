package sitecontact

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ucSiteContact "github.com/tapiaw38/practiq-be/internal/usecases/site_contact"
)

func NewGetHandler(uc ucSiteContact.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), gin.H{"message": appErr.Message()})
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
