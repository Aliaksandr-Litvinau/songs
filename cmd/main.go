package main

import (
	"fmt"
	"gin/internal/app/api"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// _ "gin/docs"
	"gin/internal/app/config"
	//"gin/internal/app/models"
	pg "gin/internal/pkg"
	//"github.com/swaggo/files"
	//"github.com/swaggo/gin-swagger"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
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

	r := api.SetupRouter()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Println("failed to run server: %v", err)
	}

	return nil
}
