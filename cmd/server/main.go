package main

import (
	"gin/config"
	"gin/internal/api"
	"gin/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	config := config.LoadConfig()

	dsn := "host=db user=user dbname=music_library port=5432 password=password"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatal("failed to connect database: ", err)
	}

	if err := db.AutoMigrate(&models.Song{}); err != nil {
		logrus.Fatal("failed to migrate database: ", err)
	}

	api.Db = db

	r := gin.Default()

	r.GET("/songs", api.GetSongs)
	r.POST("/songs", api.AddSong)

	r.Run(":" + config.Port)
}
