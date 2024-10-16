package main

import (
	"gin/config"
	_ "gin/docs"
	"gin/internal/api"
	"gin/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	cfg := config.LoadConfig()
	dsn := cfg.DatabaseURL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logrus.Fatal("failed to connect database: ", err)
	}

	if err := db.AutoMigrate(&models.Song{}); err != nil {
		logrus.Fatal("failed to migrate database: ", err)
	}

	api.Db = db

	r := api.SetupRouter()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		logrus.Fatalf("failed to run server: %v", err)
	}
}
