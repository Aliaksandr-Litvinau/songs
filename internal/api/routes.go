package api

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/songs", GetSongs)
		api.POST("/songs", AddSong)
	}

	return r
}
