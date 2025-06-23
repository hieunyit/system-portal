// internal/shared/response/router.go
package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "system-portal/docs" // Import generated Swagger docs
	authRoutes "system-portal/internal/domains/auth/routes"
	openvpnRoutes "system-portal/internal/domains/openvpn/routes"
	portalRoutes "system-portal/internal/domains/portal/routes"
	"system-portal/internal/shared/middleware"
	response "system-portal/internal/shared/response"
	"system-portal/pkg/logger"
)

type RouterConfig struct {
	Port            string // ✅ Add missing Port
	Mode            string
	TimeoutDuration time.Duration
	ReadTimeout     time.Duration // ✅ Add server timeouts
	WriteTimeout    time.Duration // ✅ Add server timeouts
}

type Router struct {
	config               *RouterConfig
	authMiddleware       *middleware.AuthMiddleware
	corsMiddleware       *middleware.CorsMiddleware
	validationMiddleware *middleware.ValidationMiddleware
	auditMiddleware      *middleware.AuditMiddleware
}

func NewRouter(
	config *RouterConfig,
	authMiddleware *middleware.AuthMiddleware,
	corsMiddleware *middleware.CorsMiddleware,
	validationMiddleware *middleware.ValidationMiddleware,
	auditMiddleware *middleware.AuditMiddleware,
) *Router {
	return &Router{
		config:               config,
		authMiddleware:       authMiddleware,
		corsMiddleware:       corsMiddleware,
		validationMiddleware: validationMiddleware,
		auditMiddleware:      auditMiddleware,
	}
}

// SetupRoutes creates and configures the Gin router
func (r *Router) SetupRoutes() *gin.Engine {
	// Set Gin mode
	gin.SetMode(r.config.Mode)

	router := gin.New()

	// Disable automatic redirect for trailing slash
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(r.corsMiddleware.Handler())
	router.Use(r.corsMiddleware.SecurityHeaders())
	router.Use(r.validationMiddleware.StrictJSONBinding())
	if r.auditMiddleware != nil {
		router.Use(r.auditMiddleware.Handler())
	}

	// Timeout middleware
	router.Use(timeout.New(
		timeout.WithTimeout(r.config.TimeoutDuration),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
	))

	// Health check and API info
	r.setupSystemRoutes(router)

	// Public routes (no authentication required)
	r.setupPublicRoutes(router)

	// Protected routes (authentication required)
	r.setupProtectedRoutes(router)

	return router
}

// ✅ NEW: StartServer starts the HTTP server with graceful shutdown
func (r *Router) StartServer() error {
	router := r.SetupRoutes()

	// Create HTTP server with timeouts
	server := &http.Server{
		Addr:         ":" + r.config.Port,
		Handler:      router,
		ReadTimeout:  r.config.ReadTimeout,
		WriteTimeout: r.config.WriteTimeout,
		IdleTimeout:  120 * time.Second, // Connection keep-alive timeout
	}

	// Start server in goroutine
	go func() {
		r.logStartupInfo()

		logger.Log.Infof("Starting System Portal API on port %s", r.config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	return r.waitForShutdown(server)
}

// ✅ NEW: StartServerWithTLS starts HTTPS server (for production)
func (r *Router) StartServerWithTLS(certFile, keyFile string) error {
	router := r.SetupRoutes()

	server := &http.Server{
		Addr:         ":" + r.config.Port,
		Handler:      router,
		ReadTimeout:  r.config.ReadTimeout,
		WriteTimeout: r.config.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		r.logStartupInfo()

		logger.Log.Infof("Starting HTTPS System Portal API on port %s", r.config.Port)
		if err := server.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("Failed to start HTTPS server: %v", err)
		}
	}()

	return r.waitForShutdown(server)
}

// ✅ NEW: logStartupInfo logs server startup information
func (r *Router) logStartupInfo() {
	logger.Log.Info(strings.Repeat("=", 60))
	logger.Log.Info("Service: System Portal API")
	logger.Log.Info("Version: 2.0.0")
	logger.Log.Infof("Server: http://localhost:%s", r.config.Port)
	logger.Log.Infof("Documentation: http://localhost:%s/swagger/index.html", r.config.Port)
	logger.Log.Infof("Health Check: http://localhost:%s/health", r.config.Port)
	logger.Log.Infof("Mode: %s", r.config.Mode)
	logger.Log.Infof("Request Timeout: %v", r.config.TimeoutDuration)
	logger.Log.Infof("Read Timeout: %v", r.config.ReadTimeout)
	logger.Log.Infof("Write Timeout: %v", r.config.WriteTimeout)
	logger.Log.Info(strings.Repeat("=", 60))
	logger.Log.Info("Domains Available: auth, portal, openvpn")
	logger.Log.Info(strings.Repeat("=", 60))
}

// ✅ NEW: waitForShutdown handles graceful shutdown
func (r *Router) waitForShutdown(server *http.Server) error {
	// Create channel to receive OS signals
	quit := make(chan os.Signal, 1)

	// Register channel to receive specific signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal is received
	sig := <-quit
	logger.Log.Infof("Received signal: %v. Shutting down server...", sig)

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Disable keep-alives and shutdown gracefully
	server.SetKeepAlivesEnabled(false)

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Errorf("Server forced to shutdown: %v", err)
		return err
	}

	logger.Log.Info("Server exited gracefully")
	return nil
}

func (r *Router) setupSystemRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", r.healthCheck)

	// API information
	router.GET("/", r.apiInfo)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (r *Router) setupPublicRoutes(router *gin.Engine) {
	// Auth routes (public - no authentication required)
	authRoutes.RegisterPublicRoutes(router)
}

func (r *Router) setupProtectedRoutes(router *gin.Engine) {
	// Protected routes group
	protected := router.Group("/")
	protected.Use(r.authMiddleware.RequireAuth())

	// Register domain routes
	authRoutes.RegisterProtectedRoutes(protected)
	portalRoutes.RegisterRoutes(protected)
	openvpnRoutes.SetRouterGroup(protected)
}

// ✅ ENHANCED: healthCheck with more detailed information
func (r *Router) healthCheck(c *gin.Context) {
	// Basic health check
	health := gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "system-portal-api",
		"version":   "2.0.0",
		"uptime":    time.Since(startTime).String(), // You'll need to track start time
	}

	// Add system information
	health["system"] = gin.H{
		"domains": []string{"auth", "portal", "openvpn"},
		"features": []string{
			"domain-architecture",
			"postgresql-auth",
			"rbac-permissions",
			"audit-logging",
			"redis-caching",
			"bulk-operations",
			"advanced-search",
		},
		"config": gin.H{
			"mode":          r.config.Mode,
			"timeout":       r.config.TimeoutDuration.String(),
			"read_timeout":  r.config.ReadTimeout.String(),
			"write_timeout": r.config.WriteTimeout.String(),
		},
	}

	response.RespondWithSuccess(c, 200, health)
}

func (r *Router) apiInfo(c *gin.Context) {
	response.RespondWithSuccess(c, 200, gin.H{
		"service":     "System Portal API",
		"version":     "2.0.0",
		"description": "Domain-driven OpenVPN Access Server Management API with PostgreSQL Authentication",
		"architecture": gin.H{
			"pattern": "Domain-Driven Design",
			"domains": []string{"auth", "portal", "openvpn"},
			"cache":   "Redis with domain namespacing",
		},
		"features": gin.H{
			"authentication": gin.H{
				"type":        "PostgreSQL + JWT",
				"permissions": "RBAC with groups",
				"audit":       "Comprehensive logging",
			},
			"openvpn_management": gin.H{
				"users":  "Full CRUD with permissions",
				"groups": "Management with access control",
				"bulk":   "Import/export operations",
				"search": "Advanced filtering",
				"cache":  "Redis performance optimization",
			},
			"portal_management": gin.H{
				"users":     "Portal user administration",
				"groups":    "Role-based group management",
				"audit":     "Activity logging and reporting",
				"dashboard": "Statistics and monitoring",
			},
		},
		"endpoints": gin.H{
			"auth":    "/auth/*",
			"portal":  "/api/portal/*",
			"openvpn": "/api/openvpn/*",
			"docs":    "/swagger/index.html",
		},
		"documentation": gin.H{
			"swagger_ui":   "/swagger/index.html",
			"swagger_json": "/swagger/doc.json",
			"api_info":     "/",
			"health":       "/health",
		},
	})
}

// ✅ NEW: Global variable to track start time (for uptime calculation)
var startTime = time.Now()
