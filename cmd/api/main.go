// cmd/api/main.go
package main

import (
	"log"
	"time"

	authHandlers "system-portal/internal/domains/auth/handlers"
	authRoutes "system-portal/internal/domains/auth/routes"
	authUsecases "system-portal/internal/domains/auth/usecases"
	openvpnHandlers "system-portal/internal/domains/openvpn/handlers"
	openvpnRepo "system-portal/internal/domains/openvpn/repositories/impl"
	openvpnRoutes "system-portal/internal/domains/openvpn/routes"
	openvpnUsecases "system-portal/internal/domains/openvpn/usecases"
	portalHandlers "system-portal/internal/domains/portal/handlers"
	portalRepo "system-portal/internal/domains/portal/repositories/impl"
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
	if err := db.Migrate(); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	// Initialize JWT service
	jwtService, err := jwt.NewRSAService(cfg.JWT.AccessTokenExpireDuration, cfg.JWT.RefreshTokenExpireDuration)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	corsMiddleware := middleware.NewCorsMiddleware(cfg.Security.CORS)
	validationMiddleware := middleware.NewValidationMiddleware()

	// Initialize domain handlers and routes
	initializeDomainRoutes(cfg, db, jwtService)

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

func initializeDomainRoutes(cfg *config.Config, db *database.Postgres, jwtSvc *jwt.RSAService) {
	// Auth domain
	authUsecase := authUsecases.NewAuthUsecase(jwtSvc)
	authHandler := authHandlers.NewAuthHandler(authUsecase)
	authRoutes.Initialize(authHandler)

	// Portal domain using PostgreSQL repositories
	userRepo := portalRepo.NewUserRepositoryPG(db.DB)
	groupRepo := portalRepo.NewGroupRepositoryPG(db.DB)
	auditRepo := portalRepo.NewAuditRepositoryPG(db.DB)

	userUC := portalUsecases.NewUserUsecase(userRepo)
	groupUC := portalUsecases.NewGroupUsecase(groupRepo)
	auditUC := portalUsecases.NewAuditUsecase(auditRepo)

	userHandler := portalHandlers.NewUserHandler(userUC)
	groupHandler := portalHandlers.NewGroupHandler(groupUC)
	auditHandler := portalHandlers.NewAuditHandler(auditUC)
	dashboardHandler := portalHandlers.NewDashboardHandler(userRepo, auditRepo)

	portalRoutes.Initialize(userHandler, groupHandler, auditHandler, dashboardHandler)

	// OpenVPN domain initialization
	xmlrpcClient := xmlrpc.NewClient(cfg.OpenVPN)
	ldapClient := ldap.NewClient(cfg.LDAP)

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
