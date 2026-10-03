package upload

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-be/internal/platform/storage"
	ucUpload "github.com/tapiaw38/practiq-be/internal/usecases/upload"
)

var allowedFolders = map[string]bool{
	"attachments": true,
	"materials":   true,
	"exercises":   true,
}

const multipartOverhead = 1 << 20

func NewHandler(uc ucUpload.UploadUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, storage.MaxUploadBytes+multipartOverhead)

		fileHeader, err := c.FormFile("file")
		if err != nil {
			var maxBytesErr *http.MaxBytesError
			if errors.As(err, &maxBytesErr) {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{"code": "upload:too-large", "message": "file is too large"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "a file field is required"})
			return
		}

		folder := c.PostForm("folder")
		if folder == "" {
			folder = "attachments"
		}
		if !allowedFolders[folder] {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "invalid folder"})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "could not read the file"})
			return
		}
		defer file.Close()

		output, appErr := uc.Execute(c, ucUpload.UploadInput{
			UserID:      middlewares.GetUserID(c),
			Folder:      folder,
			Filename:    fileHeader.Filename,
			ContentType: fileHeader.Header.Get("Content-Type"),
			Reader:      file,
			Size:        fileHeader.Size,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
