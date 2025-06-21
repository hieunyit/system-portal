// cmd/api/main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	authHandlers "system-portal/internal/domains/auth/handlers"
	sessionRepoimpl "system-portal/internal/domains/auth/repositories/impl"
	authRoutes "system-portal/internal/domains/auth/routes"
	authUsecases "system-portal/internal/domains/auth/usecases"
	openvpnHandlers "system-portal/internal/domains/openvpn/handlers"
	openvpnRepo "system-portal/internal/domains/openvpn/repositories/impl"
	openvpnRoutes "system-portal/internal/domains/openvpn/routes"
	openvpnUsecases "system-portal/internal/domains/openvpn/usecases"
	"system-portal/internal/domains/portal/entities"
	portalHandlers "system-portal/internal/domains/portal/handlers"
	portalRepo "system-portal/internal/domains/portal/repositories"
	portalRepoImpl "system-portal/internal/domains/portal/repositories/impl"
	portalRoutes "system-portal/internal/domains/portal/routes"
	portalUsecases "system-portal/internal/domains/portal/usecases"
	"system-portal/internal/shared/config"
	"system-portal/internal/shared/database"
	serverHttp "system-portal/internal/shared/infrastructure/http"
	"system-portal/internal/shared/infrastructure/ldap"
	"system-portal/internal/shared/infrastructure/xmlrpc"
	"system-portal/internal/shared/middleware"
	"system-portal/pkg/jwt"
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

	// Connect to PostgreSQL and run migrations
	db, err := database.New(cfg.Database.DSN)
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	defer db.Close()

	logger.Log.WithField("dsn", cfg.Database.DSN).Info("checking database connectivity")
	if err := waitForPostgres(db.DB, 5, time.Second); err != nil {
		log.Fatal("database unreachable:", err)
	}

	logDatabaseVersion(db.DB)

	if err := db.Migrate(); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	logExistingUsers(db.DB)

	// Initialize JWT service
	jwtService, err := jwt.NewRSAService(cfg.JWT.AccessTokenExpireDuration, cfg.JWT.RefreshTokenExpireDuration)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize infrastructure clients
	xmlrpcClient := xmlrpc.NewClient(cfg.OpenVPN)
	ldapClient := ldap.NewClient(cfg.LDAP)

	// Verify external service connections
	if err := checkConnections(db, ldapClient, xmlrpcClient); err != nil {
		log.Fatal("connectivity check failed:", err)
	}

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	corsMiddleware := middleware.NewCorsMiddleware(cfg.Security.CORS)
	validationMiddleware := middleware.NewValidationMiddleware()

	// Initialize domain handlers and routes
	initializeDomainRoutes(cfg, db, jwtService, xmlrpcClient, ldapClient)

	// Create router configuration
	routerConfig := &serverHttp.RouterConfig{
		Port:            cfg.Server.Port,
		Mode:            cfg.Server.Mode,
		TimeoutDuration: time.Duration(cfg.Server.Timeout) * time.Second,
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
	}

	// Create router
	router := serverHttp.NewRouter(
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

func initializeDomainRoutes(cfg *config.Config, db *database.Postgres, jwtSvc *jwt.RSAService, xmlrpcClient *xmlrpc.Client, ldapClient *ldap.Client) {
	// Portal domain using PostgreSQL repositories
	userRepo := portalRepoImpl.NewUserRepositoryPG(db.DB)
	groupRepo := portalRepoImpl.NewGroupRepositoryPG(db.DB)
	auditRepo := portalRepoImpl.NewAuditRepositoryPG(db.DB)

	// Auth domain
	sessionRepo := sessionRepoimpl.NewSessionRepositoryPG(db.DB)
	authUsecase := authUsecases.NewAuthUsecase(sessionRepo, userRepo, jwtSvc)
	authHandler := authHandlers.NewAuthHandler(authUsecase)
	authRoutes.Initialize(authHandler)

	userUC := portalUsecases.NewUserUsecase(userRepo)
	groupUC := portalUsecases.NewGroupUsecase(groupRepo)
	auditUC := portalUsecases.NewAuditUsecase(auditRepo)

	userHandler := portalHandlers.NewUserHandler(userUC)
	groupHandler := portalHandlers.NewGroupHandler(groupUC)
	auditHandler := portalHandlers.NewAuditHandler(auditUC)
	dashboardHandler := portalHandlers.NewDashboardHandler(userRepo, auditRepo)

	portalRoutes.Initialize(userHandler, groupHandler, auditHandler, dashboardHandler)

	// OpenVPN domain initialization

	userRepoOV := openvpnRepo.NewUserRepository(xmlrpcClient)
	groupRepoOV := openvpnRepo.NewGroupRepository(xmlrpcClient)
	disconnectRepo := openvpnRepo.NewDisconnectRepository(xmlrpcClient)
	vpnStatusRepo := openvpnRepo.NewVPNStatusRepository(xmlrpcClient)
	configRepoOV := openvpnRepo.NewConfigRepository(xmlrpcClient)

	userUCOV := openvpnUsecases.NewUserUsecase(userRepoOV, groupRepoOV, ldapClient)
	groupUCOV := openvpnUsecases.NewGroupUsecase(groupRepoOV, configRepoOV)
	bulkUCOV := openvpnUsecases.NewBulkUsecase(userRepoOV, groupRepoOV, ldapClient)
	disconnectUC := openvpnUsecases.NewDisconnectUsecase(userRepoOV, disconnectRepo, vpnStatusRepo)
	configUCOV := openvpnUsecases.NewConfigUsecase(configRepoOV)
	vpnStatusUC := openvpnUsecases.NewVPNStatusUsecase(vpnStatusRepo)

	userHandlerOV := openvpnHandlers.NewUserHandler(userUCOV, xmlrpcClient)
	groupHandlerOV := openvpnHandlers.NewGroupHandler(groupUCOV, configUCOV, xmlrpcClient)
	bulkHandlerOV := openvpnHandlers.NewBulkHandler(bulkUCOV, xmlrpcClient)
	configHandlerOV := openvpnHandlers.NewConfigHandler(configUCOV)
	vpnStatusHandlerOV := openvpnHandlers.NewVPNStatusHandler(vpnStatusUC)
	disconnectHandlerOV := openvpnHandlers.NewDisconnectHandler(disconnectUC)

	openvpnRoutes.Initialize(
		userHandlerOV,
		groupHandlerOV,
		bulkHandlerOV,
		configHandlerOV,
		vpnStatusHandlerOV,
		disconnectHandlerOV,
	)
}

// waitForPostgres pings the database until it responds or retries are exhausted.
func waitForPostgres(db *sql.DB, retries int, delay time.Duration) error {
	for i := 0; i < retries; i++ {
		if err := db.Ping(); err == nil {
			logger.Log.Info("postgres connection established")
			return nil
		} else {
			logger.Log.WithError(err).Warnf("postgres ping failed, retry %d/%d", i+1, retries)
			time.Sleep(delay)
		}
	}
	return fmt.Errorf("unable to reach postgres after %d attempts", retries)
}

// checkConnections verifies connectivity to PostgreSQL, LDAP and OpenVPN XML-RPC services.
func checkConnections(db *database.Postgres, ldapClient *ldap.Client, xmlClient *xmlrpc.Client) error {
	// Verify PostgreSQL connection
	if err := db.DB.Ping(); err != nil {
		return fmt.Errorf("postgres connection failed: %w", err)
	}

	// Verify LDAP connection
	conn, err := ldapClient.Connect()
	if err != nil {
		return fmt.Errorf("ldap connection failed: %w", err)
	}
	conn.Close()

	// Verify OpenVPN XML-RPC endpoint
	if err := xmlClient.Ping(); err != nil {
		return fmt.Errorf("openvpn connection failed: %w", err)
	}

	return nil
}

// logExistingUsers outputs all records from the users table for debugging.
func logExistingUsers(db *sql.DB) {
	repo := portalRepoImpl.NewUserRepositoryPG(db)
	users, err := repo.List(context.Background())
	if err != nil {
		logger.Log.WithError(err).Error("failed to list users")
		return
	}
	logger.Log.WithField("total", len(users)).Info("users table loaded")
	for _, u := range users {
		logger.Log.WithFields(map[string]interface{}{
			"id":       u.ID,
			"username": u.Username,
			"email":    u.Email,
			"active":   u.IsActive,
		}).Info("db user")
	}
}

// logDatabaseVersion queries and logs the PostgreSQL server version.
func logDatabaseVersion(db *sql.DB) {
	var version string
	if err := db.QueryRow(`SELECT version()`).Scan(&version); err != nil {
		logger.Log.WithError(err).Warn("failed to retrieve postgres version")
		return
	}
	logger.Log.WithField("version", version).Info("postgres version")
}
