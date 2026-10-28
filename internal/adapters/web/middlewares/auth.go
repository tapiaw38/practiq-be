package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	userprofile "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/user_profile"
	"github.com/tapiaw38/practiq-be/internal/platform/auth"
	"github.com/tapiaw38/practiq-be/internal/platform/revocation"
)

func AuthMiddleware(checker *revocation.Checker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "common:unauthorized", "message": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "common:unauthorized", "message": "invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "common:unauthorized", "message": "invalid token"})
			c.Abort()
			return
		}

		if checker != nil && !checker.Current(c.Request.Context(), authHeader, claims.UserID, claims.TokenVersion) {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "common:unauthorized", "message": "session ended, please sign in again"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

func GetUserID(c *gin.Context) string {
	userID, _ := c.Get("userID")
	if id, ok := userID.(string); ok {
		return id
	}
	return ""
}

func GetRoles(c *gin.Context) []auth.RoleClaim {
	roles, _ := c.Get("roles")
	if value, ok := roles.([]auth.RoleClaim); ok {
		return value
	}
	return nil
}

func HasRole(c *gin.Context, expected ...string) bool {
	roles := GetRoles(c)
	for _, role := range roles {
		for _, exp := range expected {
			if role.Name == exp {
				return true
			}
		}
	}
	return false
}

const (
	RoleSuperAdmin = "superadmin"
)

func IsSuperAdmin(c *gin.Context) bool {
	return HasRole(c, RoleSuperAdmin)
}

func IsTeacher(c *gin.Context) bool {
	if IsSuperAdmin(c) {
		return true
	}
	isTeacher, _ := c.Get("isTeacher")
	value, _ := isTeacher.(bool)
	return value
}

func LoadProfileType(profiles userprofile.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		profile, err := profiles.Get(context.Background(), GetUserID(c))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "common:internal-error", "message": "failed to load profile"})
			c.Abort()
			return
		}
		c.Set("isTeacher", profile != nil && profile.ProfileType == "teacher")
		c.Next()
	}
}

func RequireTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsTeacher(c) {
			c.JSON(http.StatusForbidden, gin.H{"code": "common:forbidden", "message": "teacher profile required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireRoles(expected ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !HasRole(c, expected...) {
			c.JSON(http.StatusForbidden, gin.H{"code": "common:forbidden", "message": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}
