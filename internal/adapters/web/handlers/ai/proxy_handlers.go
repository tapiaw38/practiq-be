package ai

import (
	"github.com/gin-gonic/gin"
)

func messagePath(c *gin.Context) string {
	path := "/conversation/" + c.Param("id") + "/message"
	if rawQuery := c.Request.URL.RawQuery; rawQuery != "" {
		path += "?" + rawQuery
	}
	return path
}

func streamMessagePath(c *gin.Context) string {
	path := "/conversation/" + c.Param("id") + "/message/stream"
	if rawQuery := c.Request.URL.RawQuery; rawQuery != "" {
		path += "?" + rawQuery
	}
	return path
}
