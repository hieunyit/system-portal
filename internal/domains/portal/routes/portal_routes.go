// internal/domains/portal/routes/portal_routes.go
package routes

import (
	"system-portal/internal/domains/portal/handlers"
	"system-portal/internal/shared/middleware"

	"github.com/gin-gonic/gin"
)

// Dependencies injected from main
var (
	userHandler      *handlers.UserHandler
	groupHandler     *handlers.GroupHandler
	auditHandler     *handlers.AuditHandler
	dashboardHandler *handlers.DashboardHandler
)

// Initialize sets up the handler dependencies
func Initialize(
	uh *handlers.UserHandler,
	gh *handlers.GroupHandler,
	ah *handlers.AuditHandler,
	dh *handlers.DashboardHandler,
) {
	userHandler = uh
	groupHandler = gh
	auditHandler = ah
	dashboardHandler = dh
}

// RegisterRoutes registers all portal routes
func RegisterRoutes(router *gin.RouterGroup) {
	portal := router.Group("/api/portal")

	// All portal routes require admin access
	portal.Use(middleware.RequireGroup("admin"))

	// Portal user management routes
	registerUserRoutes(portal)

	// Portal group management routes
	registerGroupRoutes(portal)

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
		// Groups are predefined (admin, support), so no create/update/delete
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
