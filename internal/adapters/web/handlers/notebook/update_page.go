package notebook

import (
	"net/http"

	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewUpdatePageHandler(uc ucNB.UpdatePageUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageID := c.Param("id")
		requesterID := middlewares.GetUserID(c)
		isAdmin := middlewares.HasRole(c, "admin", "superadmin")
		var input struct {
			Title        string `json:"title"`
			ContentType  string `json:"content_type"`
			ContentData  string `json:"content_data"`
			Instructions string `json:"instructions"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "internal server error"})
			return
		}
		if appErr := uc.Execute(c, requesterID, isAdmin, ucNB.UpdatePageInput{
			PageID:       pageID,
			Title:        input.Title,
			ContentType:  input.ContentType,
			ContentData:  input.ContentData,
			Instructions: input.Instructions,
		}); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
