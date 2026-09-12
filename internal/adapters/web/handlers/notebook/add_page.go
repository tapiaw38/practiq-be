package notebook

import (
	"net/http"

	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewAddPageHandler(uc ucNB.AddPageUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		notebookID := c.Param("id")
		requesterID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)
		var input struct {
			PageNumber   int    `json:"page_number"`
			Title        string `json:"title"`
			ContentType  string `json:"content_type"`
			ContentData  string `json:"content_data"`
			Instructions string `json:"instructions"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "internal server error"})
			return
		}
		out, appErr := uc.Execute(c, requesterID, isSuperAdmin, ucNB.AddPageInput{
			NotebookID:   notebookID,
			PageNumber:   input.PageNumber,
			Title:        input.Title,
			ContentType:  input.ContentType,
			ContentData:  input.ContentData,
			Instructions: input.Instructions,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusCreated, out)
	}
}
