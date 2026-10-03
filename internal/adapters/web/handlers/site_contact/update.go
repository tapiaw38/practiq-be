package sitecontact

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ucSiteContact "github.com/tapiaw38/practiq-be/internal/usecases/site_contact"
)

func NewUpdateHandler(uc ucSiteContact.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucSiteContact.UpdateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), gin.H{"message": appErr.Message()})
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
