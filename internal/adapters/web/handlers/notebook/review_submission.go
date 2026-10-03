package notebook

import (
	"log"
	"net/http"

	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewReviewSubmissionHandler(uc ucNB.ReviewSubmissionUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		submissionID := c.Param("id")
		teacherID := ""
		if !middlewares.IsSuperAdmin(c) {
			teacherID = middlewares.GetUserID(c)
		}
		output, err := uc.Execute(c, submissionID, teacherID, c.GetHeader("Authorization"))
		if err != nil {
			if err.Error() == "submission not found" {
				c.JSON(http.StatusNotFound, gin.H{"code": "notebook:submission-not-found", "message": "submission not found"})
				return
			}
			log.Printf("[notebook handler] review submission error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "notebook:review-error", "message": "internal server error"})
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
