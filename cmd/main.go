package main

import (
	"fmt"
	// _ "gin/docs"
	"gin/internal/app/config"
	//"gin/internal/app/models"
	pg "gin/internal/pkg"
	//"github.com/swaggo/files"
	//"github.com/swaggo/gin-swagger"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run() error {
	// read config
	cfg := config.Read()

	pgDB, err := pg.Dial(cfg.DSN)
	if err != nil {
		return fmt.Errorf("pg.Dial failed: %w", err)
	}

	// TODO: add migrations
	if pgDB != nil {
		log.Println("Running PostgreSQL migrations")
		//if err := runPgMigrations(cfg.DSN, cfg.MigrationsPath); err != nil {
		//	return fmt.Errorf("runPgMigrations failed: %w", err)
		//}
	}

	return nil
}
