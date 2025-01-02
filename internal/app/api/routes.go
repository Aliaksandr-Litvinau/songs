package api

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	r := gin.Default()

	//api := r.Group("/api")
	//{
	//	api.GET("/songs", GetSongs)
	//	api.GET("/songs/:id", GetSong)
	//	api.POST("/songs", AddSong)
	//	api.PUT("/songs/:id", UpdateSong)
	//	api.PATCH("/songs/:id", PartialUpdateSong)
	//	api.DELETE("/songs/:id", DeleteSong)
	//	api.GET("/songs/:id/verses", GetSongVerses)
	//}

	return r
}
