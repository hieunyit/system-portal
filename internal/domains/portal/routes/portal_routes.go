// internal/domains/portal/routes/portal_routes.go
package routes

import (
	portalHandlers "system-portal/internal/domains/portal/handlers"
	"system-portal/internal/shared/middleware"

	"github.com/gin-gonic/gin"
)

// Dependencies injected from main
var (
	userHandler       *portalHandlers.UserHandler
	groupHandler      *portalHandlers.GroupHandler
	permissionHandler *portalHandlers.PermissionHandler
	auditHandler      *portalHandlers.AuditHandler
	dashboardHandler  *portalHandlers.DashboardHandler
	configHandler     *portalHandlers.ConfigHandler
)

// Initialize sets up the handler dependencies
func Initialize(
	uh *portalHandlers.UserHandler,
	gh *portalHandlers.GroupHandler,
	ph *portalHandlers.PermissionHandler,
	ah *portalHandlers.AuditHandler,
	dh *portalHandlers.DashboardHandler,
	ch *portalHandlers.ConfigHandler,
) {
	userHandler = uh
	groupHandler = gh
	permissionHandler = ph
	auditHandler = ah
	dashboardHandler = dh
	configHandler = ch
}

// RegisterRoutes registers all portal routes
func RegisterRoutes(router *gin.RouterGroup) {
	portal := router.Group("/api/portal")

	// All portal routes require admin access
	portal.Use(middleware.RequireGroup("admin"))

	// Portal user management routes
	registerUserRoutes(portal)
	registerPermissionRoutes(portal)

	// Portal group management routes
	registerGroupRoutes(portal)

	// Connection config routes
	registerConfigRoutes(portal)

	// Audit log routes
	registerAuditRoutes(portal)

	// Dashboard routes
	registerDashboardRoutes(portal)
}

func registerUserRoutes(portal *gin.RouterGroup) {
	users := portal.Group("/users")
	{
		users.GET("", userHandler.ListUsers)
		users.POST("", userHandler.CreateUser)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)

		// User actions
		users.PUT("/:id/activate", userHandler.ActivateUser)
		users.PUT("/:id/deactivate", userHandler.DeactivateUser)
		users.PUT("/:id/reset-password", userHandler.ResetPassword)
	}
}

func registerGroupRoutes(portal *gin.RouterGroup) {
	groups := portal.Group("/groups")
	{
		groups.GET("", groupHandler.ListGroups)
		groups.GET("/:id", groupHandler.GetGroup)
		groups.POST("", groupHandler.CreateGroup)
		groups.PUT("/:id", groupHandler.UpdateGroup)
		groups.DELETE("/:id", groupHandler.DeleteGroup)
		groups.GET("/:id/permissions", groupHandler.GetGroupPermissions)
		groups.PUT("/:id/permissions", groupHandler.UpdateGroupPermissions)
	}
}

func registerPermissionRoutes(portal *gin.RouterGroup) {
	perms := portal.Group("/permissions")
	{
		perms.GET("", permissionHandler.ListPermissions)
		perms.POST("", permissionHandler.CreatePermission)
		perms.PUT("/:id", permissionHandler.UpdatePermission)
		perms.DELETE("/:id", permissionHandler.DeletePermission)
	}
}

func registerAuditRoutes(portal *gin.RouterGroup) {
	audit := portal.Group("/audit")
	{
		audit.GET("/logs", auditHandler.GetAuditLogs)
		audit.GET("/logs/export", auditHandler.ExportAuditLogs)
		audit.GET("/stats", auditHandler.GetAuditStats)
	}
}

func registerDashboardRoutes(portal *gin.RouterGroup) {
	dashboard := portal.Group("/dashboard")
	{
		dashboard.GET("/stats", dashboardHandler.GetDashboardStats)
		dashboard.GET("/activities", dashboardHandler.GetRecentActivities)
		dashboard.GET("/charts/users", dashboardHandler.GetUserChartData)
		dashboard.GET("/charts/activities", dashboardHandler.GetActivityChartData)
	}
}

func registerConfigRoutes(portal *gin.RouterGroup) {
	conn := portal.Group("/connections")
	{
		conn.GET("/openvpn", configHandler.GetOpenVPNConfig)
		conn.POST("/openvpn", configHandler.CreateOpenVPNConfig)
		conn.PUT("/openvpn", configHandler.UpdateOpenVPNConfig)
		conn.DELETE("/openvpn", configHandler.DeleteOpenVPNConfig)
		conn.POST("/openvpn/test", configHandler.TestOpenVPN)
		conn.GET("/ldap", configHandler.GetLDAPConfig)
		conn.POST("/ldap", configHandler.CreateLDAPConfig)
		conn.PUT("/ldap", configHandler.UpdateLDAPConfig)
		conn.DELETE("/ldap", configHandler.DeleteLDAPConfig)
		conn.POST("/ldap/test", configHandler.TestLDAP)
	}
}
