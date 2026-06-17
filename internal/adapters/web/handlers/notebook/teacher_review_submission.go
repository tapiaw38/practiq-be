package notebook

import (
	"log"
	"net/http"

	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"

	"github.com/gin-gonic/gin"
)

func NewTeacherReviewSubmissionHandler(uc ucNB.TeacherReviewSubmissionUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		submissionID := c.Param("id")
		var input struct {
			TeacherIsCorrect bool   `json:"teacher_is_correct"`
			TeacherFeedback  string `json:"teacher_feedback"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		output, err := uc.Execute(c, submissionID, ucNB.TeacherReviewInput{
			IsCorrect: input.TeacherIsCorrect,
			Feedback:  input.TeacherFeedback,
			TeacherID: teacherIDForReview(c),
		})
		if err != nil {
			if err.Error() == "submission not found" {
				c.JSON(http.StatusNotFound, gin.H{"code": "notebook:submission-not-found", "message": "submission not found"})
				return
			}
			log.Printf("[notebook handler] teacher review error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "notebook:teacher-review-error", "message": "internal server error"})
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
