package main

// @title       System Portal API
// @version     1.0
// @description API quản lý OpenVPN, Authentication, Portal...
// @host      localhost:8080
import (
	"context"
	"database/sql"
	"fmt"
	"time"

	authHandlers "system-portal/internal/domains/auth/handlers"
	sessionRepoimpl "system-portal/internal/domains/auth/repositories/impl"
	authRoutes "system-portal/internal/domains/auth/routes"
	authUsecases "system-portal/internal/domains/auth/usecases"
	openvpnHandlers "system-portal/internal/domains/openvpn/handlers"
	openvpnRepo "system-portal/internal/domains/openvpn/repositories/impl"
	openvpnRoutes "system-portal/internal/domains/openvpn/routes"
	openvpnUsecases "system-portal/internal/domains/openvpn/usecases"
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
		logger.Log.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	logger.Init(cfg.Logger)
	logger.Log.Info("Initializing System Portal...")

	// Connect to PostgreSQL and run migrations
	db, err := database.New(cfg.Database)
	if err != nil {
		logger.Log.Fatal("failed to connect database:", err)
	}
	defer db.Close()

	logger.Log.WithFields(map[string]interface{}{
		"host": cfg.Database.Host,
		"port": cfg.Database.Port,
		"db":   cfg.Database.Name,
	}).Info("checking database connectivity")
	if err := waitForPostgres(db.DB, 5, time.Second); err != nil {
		logger.Log.Fatal("database unreachable:", err)
	}

	if err := db.Migrate(); err != nil {
		logger.Log.Fatal("failed to migrate database:", err)
	}

	// Initialize JWT service
	var jwtService *jwt.RSAService
	if cfg.JWT.AccessPrivateKey != "" && cfg.JWT.RefreshPrivateKey != "" {
		jwtService, err = jwt.NewRSAServiceWithKeys(
			cfg.JWT.AccessPrivateKey,
			cfg.JWT.RefreshPrivateKey,
			cfg.JWT.AccessTokenExpireDuration,
			cfg.JWT.RefreshTokenExpireDuration,
		)
		if err != nil {
			logger.Log.Fatal("failed to load RSA keys:", err)
		}
	} else {
		jwtService, err = jwt.NewRSAService(cfg.JWT.AccessTokenExpireDuration, cfg.JWT.RefreshTokenExpireDuration)
		if err != nil {
			logger.Log.Fatal(err)
		}
		accessPEM, _ := jwtService.GetAccessPrivateKeyPEM()
		refreshPEM, _ := jwtService.GetRefreshPrivateKeyPEM()
		logger.Log.Warn("generated new RSA keys; store them in config to preserve sessions")
		logger.Log.Debug("accessPrivateKey:\n" + accessPEM)
		logger.Log.Debug("refreshPrivateKey:\n" + refreshPEM)
	}

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	corsMiddleware := middleware.NewCorsMiddleware(cfg.Security.CORS)
	validationMiddleware := middleware.NewValidationMiddleware()

	// Initialize domain handlers and routes
	auditUC, userRepo, groupRepo := initializeDomainRoutes(cfg, db, jwtService)

	auditMiddleware := middleware.NewAuditMiddleware(auditUC, userRepo, groupRepo)

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
		auditMiddleware,
	)

	// Start server with graceful shutdown
	logger.Log.Info("Starting server...")
	if err := router.StartServer(); err != nil {
		logger.Log.Fatal("Server error:", err)
	}
}

func initializeDomainRoutes(cfg *config.Config, db *database.Postgres, jwtSvc *jwt.RSAService) (portalUsecases.AuditUsecase, portalRepo.UserRepository, portalRepo.GroupRepository) {
	// Portal domain using PostgreSQL repositories
	userRepo := portalRepoImpl.NewUserRepositoryPG(db.DB)
	groupRepo := portalRepoImpl.NewGroupRepositoryPG(db.DB)
	auditRepo := portalRepoImpl.NewAuditRepositoryPG(db.DB)
	permRepo := portalRepoImpl.NewPermissionRepositoryPG(db.DB)

	// Auth domain
	sessionRepo := sessionRepoimpl.NewSessionRepositoryPG(db.DB)
	authUsecase := authUsecases.NewAuthUsecase(sessionRepo, userRepo, groupRepo, jwtSvc)
	authHandler := authHandlers.NewAuthHandler(authUsecase)
	authRoutes.Initialize(authHandler)

	userUC := portalUsecases.NewUserUsecase(userRepo, groupRepo)
	groupUC := portalUsecases.NewGroupUsecase(groupRepo, permRepo)
	permUC := portalUsecases.NewPermissionUsecase(permRepo)
	auditUC := portalUsecases.NewAuditUsecase(auditRepo)

	userHandler := portalHandlers.NewUserHandler(userUC)
	groupHandler := portalHandlers.NewGroupHandler(groupUC)
	permHandler := portalHandlers.NewPermissionHandler(permUC)
	auditHandler := portalHandlers.NewAuditHandler(auditUC)
	dashboardHandler := portalHandlers.NewDashboardHandler(userRepo, auditRepo)

	ovRepo := portalRepoImpl.NewOpenVPNConfigRepositoryPG(db.DB)
	ldapRepo := portalRepoImpl.NewLDAPConfigRepositoryPG(db.DB)
	configUC := portalUsecases.NewConfigUsecase(ovRepo, ldapRepo)
	reloadOpenVPN := configureOpenVPN(db, permRepo, groupRepo)
	configHandler := portalHandlers.NewConfigHandler(configUC, reloadOpenVPN)
	portalRoutes.Initialize(userHandler, groupHandler, permHandler, auditHandler, dashboardHandler, configHandler)

	// Initialize OpenVPN routes based on existing configs
	reloadOpenVPN()

	return auditUC, userRepo, groupRepo
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

	if ldapClient != nil {
		conn, err := ldapClient.Connect()
		if err != nil {
			return fmt.Errorf("ldap connection failed: %w", err)
		}
		conn.Close()
	}

	if xmlClient != nil {
		if err := xmlClient.Ping(); err != nil {
			return fmt.Errorf("openvpn connection failed: %w", err)
		}
	}

	return nil
}

func configureOpenVPN(db *database.Postgres, permRepo portalRepo.PermissionRepository, groupRepo portalRepo.GroupRepository) func() {
	return func() {
		ovRepo := portalRepoImpl.NewOpenVPNConfigRepositoryPG(db.DB)
		ldapRepo := portalRepoImpl.NewLDAPConfigRepositoryPG(db.DB)
		ovCfg, _ := ovRepo.Get(context.Background())
		ldapCfg, _ := ldapRepo.Get(context.Background())
		if ovCfg == nil {
			openvpnRoutes.Disable()
			return
		}
		xmlrpcClient := xmlrpc.NewClient(xmlrpc.Config{
			Host:     ovCfg.Host,
			Username: ovCfg.Username,
			Password: ovCfg.Password,
			Port:     ovCfg.Port,
		})
		var ldapClient *ldap.Client
		if ldapCfg != nil {
			ldapClient = ldap.NewClient(ldap.Config{
				Host:         ldapCfg.Host,
				Port:         ldapCfg.Port,
				BindDN:       ldapCfg.BindDN,
				BindPassword: ldapCfg.BindPassword,
				BaseDN:       ldapCfg.BaseDN,
			})
		} else {
			ldapClient = ldap.NewClient(ldap.Config{})
		}
		// Connectivity issues should not disable the API; log them for visibility
		if err := checkConnections(db, ldapClient, xmlrpcClient); err != nil {
			logger.Log.WithError(err).Warn("connectivity check failed")
		}

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
		permMiddleware := middleware.NewPermissionMiddleware(permRepo, groupRepo)

		openvpnRoutes.Initialize(
			userHandlerOV,
			groupHandlerOV,
			bulkHandlerOV,
			configHandlerOV,
			vpnStatusHandlerOV,
			disconnectHandlerOV,
			permMiddleware,
		)
	}
}
