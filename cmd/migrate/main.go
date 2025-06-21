package main

import (
	"log"

	"system-portal/internal/shared/config"
	"system-portal/internal/shared/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	pg, err := database.New(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}
	defer pg.Close()
	if err := pg.Migrate(); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("migrations applied successfully")
}
