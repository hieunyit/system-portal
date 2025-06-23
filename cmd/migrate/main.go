package main

import (
	"system-portal/internal/shared/config"
	"system-portal/internal/shared/database"
	"system-portal/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Log.Fatalf("failed to load config: %v", err)
	}
	pg, err := database.New(cfg.Database)
	if err != nil {
		logger.Log.Fatalf("database connection error: %v", err)
	}
	defer pg.Close()
	if err := pg.Migrate(); err != nil {
		logger.Log.Fatalf("migration failed: %v", err)
	}
	logger.Log.Info("migrations applied successfully")
}
