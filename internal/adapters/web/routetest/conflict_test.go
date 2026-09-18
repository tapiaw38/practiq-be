package routetest

import (
	"testing"

	"github.com/gin-gonic/gin"
)

// A static segment and a parameter share a position in the school routes:
// /schools/mine and /schools/:id. Some routers refuse that outright, and the
// refusal is a panic at registration — the server would not start at all.
func TestSchoolRoutesDoNotConflict(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("route registration panicked: %v", r)
		}
	}()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	group := router.Group("/api")
	noop := func(c *gin.Context) {}

	group.GET("/schools", noop)
	group.POST("/schools", noop)
	group.GET("/schools/mine", noop)
	group.PUT("/schools/:id", noop)
	group.POST("/schools/:id/close", noop)
	group.POST("/schools/:id/reopen", noop)
	group.GET("/schools/:id/archive", noop)
	group.GET("/schools/:id/members", noop)
	group.POST("/schools/:id/members", noop)
	group.DELETE("/schools/:id/members/:userId", noop)
}
