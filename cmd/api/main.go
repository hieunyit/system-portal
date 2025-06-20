// cmd/api/main.go
package main

import (
	"log"
	"time"

	"system-portal/internal/shared/config"
	"system-portal/internal/shared/infrastructure/http"
	"system-portal/internal/shared/middleware"
	"system-portal/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	logger.Init(cfg.Logger)
	log.Println("🚀 Initializing System Portal...")

	// Initialize infrastructure dependencies
	// ... (PostgreSQL, Redis, LDAP, XMLRPC clients setup)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	corsMiddleware := middleware.NewCorsMiddleware(cfg.Security.CORS)
	validationMiddleware := middleware.NewValidationMiddleware()

	// Initialize domain handlers and routes
	initializeDomainRoutes( /* dependencies */ )

	// Create router configuration
	routerConfig := &http.RouterConfig{
		Port:            cfg.Server.Port,
		Mode:            cfg.Server.Mode,
		TimeoutDuration: time.Duration(cfg.Server.Timeout) * time.Second,
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
	}

	// Create router
	router := http.NewRouter(
		routerConfig,
		authMiddleware,
		corsMiddleware,
		validationMiddleware,
	)

	// Start server with graceful shutdown
	log.Println("🎯 Starting server...")
	if err := router.StartServer(); err != nil {
		log.Fatal("Server error:", err)
	}
}

func initializeDomainRoutes( /* dependencies */ ) {
	// Initialize all domain routes with their dependencies
	// ... domain initialization code
}
