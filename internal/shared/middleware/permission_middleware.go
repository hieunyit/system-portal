package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequirePermission is a very small placeholder for RBAC permission checks.
// It simply blocks non-admin roles.
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if strings.ToLower(role) != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "insufficient permissions",
					"status":  http.StatusForbidden,
				},
			})
			return
		}
		c.Next()
	}
}

// RequireGroup ensures the authenticated user belongs to the specified group.
func RequireGroup(group string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if strings.ToLower(role) != strings.ToLower(group) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "forbidden",
					"status":  http.StatusForbidden,
				},
			})
			return
		}
		c.Next()
	}
}
