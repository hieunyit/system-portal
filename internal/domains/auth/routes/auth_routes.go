// internal/domains/auth/routes/auth_routes.go
package routes

import (
	"system-portal/internal/domains/auth/handlers"

	"github.com/gin-gonic/gin"
)

// Dependencies injected from main
var (
	authHandler *handlers.AuthHandler
)

// Initialize sets up the handler dependencies
func Initialize(ah *handlers.AuthHandler) {
	authHandler = ah
}

// RegisterPublicRoutes registers auth routes that don't require authentication
func RegisterPublicRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}
}

// RegisterProtectedRoutes registers auth routes that require authentication
func RegisterProtectedRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.GET("/validate", authHandler.ValidateToken)
		auth.POST("/logout", authHandler.Logout)
	}
}
