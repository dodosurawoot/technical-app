package main

import (
	"context"
	"log"

	"airclean-tracker/backend/internal/api"
	"airclean-tracker/backend/internal/auth"
	"airclean-tracker/backend/internal/config"
	"airclean-tracker/backend/internal/db"
	"airclean-tracker/backend/internal/importer"
)

func main() {
	cfg := config.Load()
	conn, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	if err := db.AutoMigrate(conn); err != nil {
		log.Fatalf("migrate database: %v", err)
	}
	if cfg.AutoImportExcel && cfg.ImportExcelPath != "" {
		result, err := importer.New(conn).ImportPath(cfg.ImportExcelPath)
		if err != nil {
			log.Printf("auto import failed: %v", err)
		} else {
			log.Printf("auto import complete: %+v", result)
		}
	}
	authMW, err := auth.New(context.Background(), cfg, conn)
	if err != nil {
		log.Fatalf("setup auth: %v", err)
	}
	server := api.New(cfg, conn, authMW)
	if err := server.Router().Run(":" + cfg.Port); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
