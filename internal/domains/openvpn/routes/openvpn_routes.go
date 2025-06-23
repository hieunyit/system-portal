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
	permMiddleware    *middleware.PermissionMiddleware
	enabled           bool
	routerGroup       *gin.RouterGroup
	routesRegistered  bool
)

// Initialize sets up the handler dependencies
func Initialize(
	uh *handlers.UserHandler,
	gh *handlers.GroupHandler,
	bh *handlers.BulkHandler,
	cfh *handlers.ConfigHandler,
	vsh *handlers.VPNStatusHandler,
	dh *handlers.DisconnectHandler,
	pmw *middleware.PermissionMiddleware,
) {
	userHandler = uh
	groupHandler = gh
	bulkHandler = bh
	configHandler = cfh
	vpnStatusHandler = vsh
	disconnectHandler = dh
	permMiddleware = pmw
	enabled = true
	if routerGroup != nil && !routesRegistered {
		RegisterRoutes(routerGroup)
		routesRegistered = true
	}
}

// Enabled reports whether OpenVPN routes are initialized
func Enabled() bool { return enabled }

// SetRouterGroup stores the router group for dynamic registration
func SetRouterGroup(rg *gin.RouterGroup) {
	routerGroup = rg
	if enabled && !routesRegistered {
		RegisterRoutes(routerGroup)
		routesRegistered = true
	}
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
		users.GET("", permMiddleware.RequirePermission("openvpn.view_users"), userHandler.ListUsers)
		users.GET("/expirations", permMiddleware.RequirePermission("openvpn.view_users"), userHandler.GetUserExpirations)
		users.GET("/:username", permMiddleware.RequirePermission("openvpn.view_users"), userHandler.GetUser)

		// Create and edit users (both admin and support)
		users.POST("", permMiddleware.RequirePermission("openvpn.create_users"), userHandler.CreateUser)
		users.PUT("/:username", permMiddleware.RequirePermission("openvpn.edit_users"), userHandler.UpdateUser)

		// User actions (both admin and support can enable/disable)
		users.PUT("/:username/:action", permMiddleware.RequirePermission("openvpn.edit_users"), userHandler.UserAction)

		// Delete users (admin only)
		users.DELETE("/:username", permMiddleware.RequirePermission("openvpn.delete_users"), userHandler.DeleteUser)

		// Disconnect users (both admin and support)
		users.POST("/:username/disconnect", permMiddleware.RequirePermission("openvpn.edit_users"), disconnectHandler.DisconnectUser)
	}
}

func registerGroupRoutes(openvpn *gin.RouterGroup) {
	groups := openvpn.Group("/groups")
	{
		// View groups (both admin and support)
		groups.GET("", permMiddleware.RequirePermission("openvpn.view_groups"), groupHandler.ListGroups)
		groups.GET("/:groupName", permMiddleware.RequirePermission("openvpn.view_groups"), groupHandler.GetGroup)

		// Manage groups (admin only)
		groups.POST("", permMiddleware.RequirePermission("openvpn.manage_groups"), groupHandler.CreateGroup)
		groups.PUT("/:groupName", permMiddleware.RequirePermission("openvpn.manage_groups"), groupHandler.UpdateGroup)
		groups.DELETE("/:groupName", permMiddleware.RequirePermission("openvpn.manage_groups"), groupHandler.DeleteGroup)
		groups.PUT("/:groupName/:action", permMiddleware.RequirePermission("openvpn.manage_groups"), groupHandler.GroupAction)
	}
}

func registerBulkRoutes(openvpn *gin.RouterGroup) {
	bulk := openvpn.Group("/bulk")
	{
		// User bulk operations
		userBulk := bulk.Group("/users")
		{
			// Create and import (both admin and support)
			userBulk.POST("/create", permMiddleware.RequirePermission("openvpn.create_users"), bulkHandler.BulkCreateUsers)
			userBulk.POST("/import", permMiddleware.RequirePermission("openvpn.create_users"), bulkHandler.ImportUsers)
			userBulk.GET("/template", permMiddleware.RequirePermission("openvpn.view_users"), bulkHandler.ExportUserTemplate)

			// Bulk actions (admin and support, but no delete for support)
			userBulk.POST("/actions", permMiddleware.RequirePermission("openvpn.edit_users"), bulkHandler.BulkUserActions)
			userBulk.POST("/extend", permMiddleware.RequirePermission("openvpn.edit_users"), bulkHandler.BulkExtendUsers)
			userBulk.POST("/disconnect", permMiddleware.RequirePermission("openvpn.edit_users"), disconnectHandler.BulkDisconnectUsers)
		}

		// Group bulk operations (admin only)
		groupBulk := bulk.Group("/groups")
		groupBulk.Use(permMiddleware.RequirePermission("openvpn.manage_groups"))
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
		config.GET("/server/info", permMiddleware.RequirePermission("openvpn.view_status"), configHandler.GetServerInfo)
		config.GET("/network", permMiddleware.RequirePermission("openvpn.view_status"), configHandler.GetNetworkConfig)
	}
}

func registerVPNStatusRoutes(openvpn *gin.RouterGroup) {
	vpn := openvpn.Group("/vpn")
	{
		// View VPN status (both admin and support)
		vpn.GET("/status", permMiddleware.RequirePermission("openvpn.view_status"), vpnStatusHandler.GetVPNStatus)
	}
}
