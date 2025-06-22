package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	portalrepos "system-portal/internal/domains/portal/repositories"
)

// PermissionMiddleware performs RBAC checks using group permissions.
type PermissionMiddleware struct {
	perms  portalrepos.PermissionRepository
	groups portalrepos.GroupRepository
}

func NewPermissionMiddleware(p portalrepos.PermissionRepository, g portalrepos.GroupRepository) *PermissionMiddleware {
	return &PermissionMiddleware{perms: p, groups: g}
}

func (m *PermissionMiddleware) RequirePermission(perm string) gin.HandlerFunc {
	parts := strings.SplitN(perm, ".", 2)
	if len(parts) != 2 {
		return func(c *gin.Context) { c.Next() }
	}
	resource, action := parts[0], parts[1]
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "forbidden", "status": http.StatusForbidden}})
			return
		}
		allowed, err := m.perms.HasGroupPermission(c.Request.Context(), role, resource, action)
		if err != nil || !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "insufficient permissions", "status": http.StatusForbidden}})
			return
		}
		c.Next()
	}
}

func RequireGroup(group string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if strings.ToLower(role) != strings.ToLower(group) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "forbidden", "status": http.StatusForbidden}})
			return
		}
		c.Next()
	}
}
