package transport

import (
	"github.com/gin-gonic/gin"
	"songs/internal/app/service"
)

func SetupRouter(svc service.SongService) *gin.Engine {
	r := gin.Default()

	handler := NewHandler(svc)

	api := r.Group("/api/v1")
	{
		api.GET("/songs", handler.GetSongs)
		api.GET("/songs/:id", handler.GetSong)
		api.POST("/songs", handler.CreateSong)
		api.PUT("/songs/:id", handler.UpdateSong)
		api.PATCH("/songs/:id", handler.PartialUpdateSong)
		api.DELETE("/songs/:id", handler.DeleteSong)
		api.GET("/songs/:id/verses", handler.GetSongVerses)
	}

	return r
}
