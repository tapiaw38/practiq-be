package practicesheet

import (
	"encoding/json"
	"log"
	"net/http"

	submitjob "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/submit_job"
	ucPS "github.com/tapiaw38/practiq-be/internal/usecases/practice_sheet"

	"github.com/gin-gonic/gin"
)

func NewGetSubmitJobHandler(repo submitjob.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		jobID := c.Param("jobId")
		job, err := repo.GetByID(c.Request.Context(), jobID)
		if err != nil {
			log.Printf("failed to get submit job: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "common:internal-error", "message": "failed to get submit job"})
			return
		}
		if job == nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "practice_sheet:submit-job-not-found", "message": "submit job not found"})
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
		if len(job.Result) > 0 {
			var result ucPS.SubmitOutput
			if err := json.Unmarshal(job.Result, &result); err == nil {
				response["result"] = &result
			}
		}
		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}
