package notebook

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	notebookRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/notebook"
	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewSaveSubmissionHandler(uc ucNB.SaveSubmissionUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageID := c.Param("id")
		studentID := middlewares.GetUserID(c)
		var input ucNB.SaveSubmissionInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "internal server error"})
			return
		}

		if err := uc.Execute(c, pageID, studentID, time.Now().UTC().UnixNano(), input); err != nil {

			if errors.Is(err, notebookRepo.ErrStaleSubmission) {
				c.JSON(http.StatusNoContent, nil)
				return
			}
			if strings.Contains(err.Error(), "forbidden") {
				c.JSON(http.StatusForbidden, gin.H{"code": "common:forbidden", "message": "forbidden"})
				return
			}
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{"code": "notebook:not-found", "message": err.Error()})
				return
			}
			log.Printf("[notebook handler] save submission error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
