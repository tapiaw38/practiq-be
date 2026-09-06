package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/school"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

const schoolHeader = "X-School-ID"

// SchoolContextMiddleware validates active-school selection. Empty selection
// remains allowed during rollout for legacy endpoints; tenant-aware handlers
// must call RequireSchool or use GetSchoolID before accessing scoped data.
func SchoolContextMiddleware(repo school.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(schoolHeader)
		if id == "" || IsSuperAdmin(c) {
			if id != "" {
				c.Set("schoolID", id)
				c.Request = c.Request.WithContext(tenant.WithSchoolID(c.Request.Context(), id))
			}
			c.Next()
			return
		}
		ok, err := repo.IsMember(c, id, GetUserID(c))
		if err != nil || !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": "common:forbidden", "message": "user is not a member of selected school"})
			return
		}
		c.Set("schoolID", id)
		c.Request = c.Request.WithContext(tenant.WithSchoolID(c.Request.Context(), id))
		c.Next()
	}
}

func GetSchoolID(c *gin.Context) string {
	id, _ := c.Get("schoolID")
	value, _ := id.(string)
	return value
}

func RequireSchool() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetSchoolID(c) == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": "common:school-required", "message": "X-School-ID header required"})
			return
		}
		c.Next()
	}
}

// RequireSchoolAdmin authorizes institutional administration from membership,
// never from auth-api's global teacher role. A teacher may teach in a school
// but cannot administer it unless that school's membership grants admin.
func RequireSchoolAdmin(repo school.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if IsSuperAdmin(c) {
			c.Next()
			return
		}
		schoolID := GetSchoolID(c)
		if schoolID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": "common:school-required", "message": "X-School-ID header required"})
			return
		}
		ok, err := repo.HasAdminAccess(c, schoolID, GetUserID(c))
		if err != nil || !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": "common:forbidden", "message": "school admin access required"})
			return
		}
		c.Next()
	}
}
