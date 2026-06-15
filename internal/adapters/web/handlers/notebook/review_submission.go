package notebook

import (
	"net/http"

	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewReviewSubmissionHandler(uc ucNB.ReviewSubmissionUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		submissionID := c.Param("id")
		teacherID := ""
		if !middlewares.HasRole(c, "admin", "superadmin") {
			teacherID = middlewares.GetUserID(c)
		}
		output, err := uc.Execute(c, submissionID, teacherID)
		if err != nil {
			if err.Error() == "submission not found" {
				c.JSON(http.StatusNotFound, gin.H{"code": "notebook:submission-not-found", "message": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": "notebook:review-error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
