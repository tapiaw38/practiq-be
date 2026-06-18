package notebook

import (
	"log"
	"net/http"

	submitjob "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/submit_job"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewGetSubmitJobHandler(repo submitjob.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		requesterID := middlewares.GetUserID(c)
		isAdmin := middlewares.HasRole(c, "admin", "superadmin")
		jobID := c.Param("jobId")
		job, err := repo.GetByID(c.Request.Context(), jobID)
		if err != nil {
			log.Printf("failed to get submit job: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "common:internal-error", "message": "failed to get submit job"})
			return
		}
		if job == nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "notebook:submit-job-not-found", "message": "submit job not found"})
			return
		}
		if !isAdmin && job.StudentID != requesterID {
			c.JSON(http.StatusForbidden, gin.H{"code": "common:forbidden", "message": "cannot view other user's submit job results"})
			return
		}
		// Preserve original JSON response shape
		response := gin.H{
			"status":     job.Status,
			"created_at": job.CreatedAt,
			"updated_at": job.UpdatedAt,
		}
		if job.ErrorCode != "" {
			response["error_code"] = job.ErrorCode
		}
		if job.Message != "" {
			response["message"] = job.Message
		}
		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}
