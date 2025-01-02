package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"songs/internal/app/api"
	"songs/internal/app/config"
	pg "songs/internal/pkg"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	if pgDB != nil {
		if err := runPgMigrations(cfg.MigrationsPath, cfg.DSN); err != nil {
			return fmt.Errorf("runPgMigrations failed: %w", err)
		}
	}

	r := api.SetupRouter()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(cfg.HTTPAddr); err != nil {
		log.Println("failed to run server: %v", err)
	}

	return nil
}

func runPgMigrations(path, dsn string) error {
	if path == "" {
		return errors.New("no migrations path provided")
	}
	if dsn == "" {
		return errors.New("no DSN provided")
	}

	log.Println("Initializing migrations")
	m, err := migrate.New(path, dsn)
	if err != nil {
		return err
	}

	log.Println("Running migrations")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations completed successfully")
	return nil
}
