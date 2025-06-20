package routes

import (
	"system-portal/internal/domains/openvpn/handlers"
	"system-portal/internal/shared/middleware"

	"github.com/gin-gonic/gin"
)

// Dependencies injected from main
var (
	userHandler  *handlers.UserHandler
	groupHandler *handlers.GroupHandler
	bulkHandler  *handlers.BulkHandler

	configHandler     *handlers.ConfigHandler
	vpnStatusHandler  *handlers.VPNStatusHandler
	disconnectHandler *handlers.DisconnectHandler
)

// Initialize sets up the handler dependencies
func Initialize(
	uh *handlers.UserHandler,
	gh *handlers.GroupHandler,
	bh *handlers.BulkHandler,
	cfh *handlers.ConfigHandler,
	vsh *handlers.VPNStatusHandler,
	dh *handlers.DisconnectHandler,
) {
	userHandler = uh
	groupHandler = gh
	bulkHandler = bh
	configHandler = cfh
	vpnStatusHandler = vsh
	disconnectHandler = dh
}

// RegisterRoutes registers all OpenVPN routes with permission-based access control
func RegisterRoutes(router *gin.RouterGroup) {
	openvpn := router.Group("/api/openvpn")

	// Register route groups
	registerUserRoutes(openvpn)
	registerGroupRoutes(openvpn)
	registerBulkRoutes(openvpn)
	registerConfigRoutes(openvpn)
	registerVPNStatusRoutes(openvpn)
}

func registerUserRoutes(openvpn *gin.RouterGroup) {
	users := openvpn.Group("/users")
	{
		// List and view users (both admin and support)
		users.GET("", middleware.RequirePermission("openvpn.view_users"), userHandler.ListUsers)
		users.GET("/expirations", middleware.RequirePermission("openvpn.view_users"), userHandler.GetUserExpirations)
		users.GET("/:username", middleware.RequirePermission("openvpn.view_users"), userHandler.GetUser)

		// Create and edit users (both admin and support)
		users.POST("", middleware.RequirePermission("openvpn.create_users"), userHandler.CreateUser)
		users.PUT("/:username", middleware.RequirePermission("openvpn.edit_users"), userHandler.UpdateUser)

		// User actions (both admin and support can enable/disable)
		users.PUT("/:username/:action", middleware.RequirePermission("openvpn.edit_users"), userHandler.UserAction)

		// Delete users (admin only)
		users.DELETE("/:username", middleware.RequirePermission("openvpn.delete_users"), userHandler.DeleteUser)

		// Disconnect users (both admin and support)
		users.POST("/:username/disconnect", middleware.RequirePermission("openvpn.edit_users"), disconnectHandler.DisconnectUser)
	}
}

func registerGroupRoutes(openvpn *gin.RouterGroup) {
	groups := openvpn.Group("/groups")
	{
		// View groups (both admin and support)
		groups.GET("", middleware.RequirePermission("openvpn.view_groups"), groupHandler.ListGroups)
		groups.GET("/:groupName", middleware.RequirePermission("openvpn.view_groups"), groupHandler.GetGroup)

		// Manage groups (admin only)
		groups.POST("", middleware.RequirePermission("openvpn.manage_groups"), groupHandler.CreateGroup)
		groups.PUT("/:groupName", middleware.RequirePermission("openvpn.manage_groups"), groupHandler.UpdateGroup)
		groups.DELETE("/:groupName", middleware.RequirePermission("openvpn.manage_groups"), groupHandler.DeleteGroup)
		groups.PUT("/:groupName/:action", middleware.RequirePermission("openvpn.manage_groups"), groupHandler.GroupAction)
	}
}

func registerBulkRoutes(openvpn *gin.RouterGroup) {
	bulk := openvpn.Group("/bulk")
	{
		// User bulk operations
		userBulk := bulk.Group("/users")
		{
			// Create and import (both admin and support)
			userBulk.POST("/create", middleware.RequirePermission("openvpn.create_users"), bulkHandler.BulkCreateUsers)
			userBulk.POST("/import", middleware.RequirePermission("openvpn.create_users"), bulkHandler.ImportUsers)
			userBulk.GET("/template", middleware.RequirePermission("openvpn.view_users"), bulkHandler.ExportUserTemplate)

			// Bulk actions (admin and support, but no delete for support)
			userBulk.POST("/actions", middleware.RequirePermission("openvpn.edit_users"), bulkHandler.BulkUserActions)
			userBulk.POST("/extend", middleware.RequirePermission("openvpn.edit_users"), bulkHandler.BulkExtendUsers)
			userBulk.POST("/disconnect", middleware.RequirePermission("openvpn.edit_users"), disconnectHandler.BulkDisconnectUsers)
		}

		// Group bulk operations (admin only)
		groupBulk := bulk.Group("/groups")
		groupBulk.Use(middleware.RequirePermission("openvpn.manage_groups"))
		{
			groupBulk.POST("/create", bulkHandler.BulkCreateGroups)
			groupBulk.POST("/actions", bulkHandler.BulkGroupActions)
			groupBulk.POST("/import", bulkHandler.ImportGroups)
			groupBulk.GET("/template", bulkHandler.ExportGroupTemplate)
		}
	}
}

func registerConfigRoutes(openvpn *gin.RouterGroup) {
	config := openvpn.Group("/config")
	{
		// View configuration (both admin and support)
		config.GET("/server/info", middleware.RequirePermission("openvpn.view_status"), configHandler.GetServerInfo)
		config.GET("/network", middleware.RequirePermission("openvpn.view_status"), configHandler.GetNetworkConfig)
	}
}

func registerVPNStatusRoutes(openvpn *gin.RouterGroup) {
	vpn := openvpn.Group("/vpn")
	{
		// View VPN status (both admin and support)
		vpn.GET("/status", middleware.RequirePermission("openvpn.view_status"), vpnStatusHandler.GetVPNStatus)
	}
}
