package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/platform/auth"
)

func AuthMiddleware() gin.HandlerFunc {
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

// Los roles vienen de auth-api-be, donde solo existen tres: user es el alumno,
// admin el profesor y superadmin el administrador de la plataforma.
const (
	RoleStudent    = "user"
	RoleTeacher    = "admin"
	RoleSuperAdmin = "superadmin"
)

// IsSuperAdmin responde si quien hace el pedido administra la plataforma. Es el
// único permiso que saltea las reglas de propiedad y de asignación
// docente-alumno, así que vive acá y no repetido en cada handler: cuando la
// condición se escribía a mano, "admin" terminó incluido en la lista y todo
// profesor quedó con el alcance del administrador.
func IsSuperAdmin(c *gin.Context) bool {
	return HasRole(c, RoleSuperAdmin)
}

// IsTeacher responde si puede dictar clase: el profesor y, por encima, el
// administrador.
func IsTeacher(c *gin.Context) bool {
	return HasRole(c, RoleTeacher, RoleSuperAdmin)
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
