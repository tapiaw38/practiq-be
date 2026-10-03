package notebook

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	notebookRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/notebook"
	submitjob "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/submit_job"
	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/utils"
)

func NewSaveSubmissionAsyncHandler(uc ucNB.SaveSubmissionUsecase, repo submitjob.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageID := c.Param("id")
		studentID := middlewares.GetUserID(c)
		var input ucNB.SaveSubmissionInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		jobID := utils.NewSubmitJobID()
		now := time.Now().UTC()

		version := now.UnixNano()
		if err := repo.Create(c.Request.Context(), domain.SubmitJob{
			ID:        jobID,
			Kind:      "notebook",
			StudentID: studentID,
			Status:    "processing",
			CreatedAt: now,
			UpdatedAt: now,
		}); err != nil {

			log.Printf("failed to create submit job: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "practice_sheet:submit-job-not-created",
				"message": "could not start the submission",
			})
			return
		}

		go func(pid, sid, jid string, ver int64, payload ucNB.SaveSubmissionInput) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			err := uc.Execute(ctx, pid, sid, ver, payload)
			finishCtx, finishCancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer finishCancel()

			if errors.Is(err, notebookRepo.ErrStaleSubmission) {
				if updateErr := repo.Update(finishCtx, domain.SubmitJob{
					ID:     jid,
					Status: "done",
				}); updateErr != nil {
					log.Printf("failed to update submit job: %v", updateErr)
				}
				return
			}
			if err != nil {
				if updateErr := repo.Update(finishCtx, domain.SubmitJob{
					ID:        jid,
					Status:    "failed",
					ErrorCode: "notebook:submit-failed",
					Message:   err.Error(),
				}); updateErr != nil {
					log.Printf("failed to update submit job: %v", updateErr)
				}
				return
			}
			if updateErr := repo.Update(finishCtx, domain.SubmitJob{
				ID:     jid,
				Status: "done",
			}); updateErr != nil {
				log.Printf("failed to update submit job: %v", updateErr)
			}
		}(pageID, studentID, jobID, version, input)

		c.JSON(http.StatusAccepted, gin.H{
			"data": gin.H{
				"job_id": jobID,
				"status": "processing",
			},
		})
	}
}
