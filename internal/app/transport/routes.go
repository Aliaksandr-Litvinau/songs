package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "songs/docs"
	"songs/internal/app/transport/adapter"
)

func SetupRouter(svc SongService) *gin.Engine {
	r := gin.Default()

	handler := NewHandler(svc)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/api/v1")
	{
		api.GET("/songs", adapter.ToGinHandler(handler.GetSongs))
		api.GET("/songs/:id", adapter.ToGinHandler(handler.GetSong))
		api.POST("/songs", adapter.ToGinHandler(handler.CreateSong))
		api.PUT("/songs/:id", adapter.ToGinHandler(handler.UpdateSong))
		api.PATCH("/songs/:id", adapter.ToGinHandler(handler.PartialUpdateSong))
		api.DELETE("/songs/:id", adapter.ToGinHandler(handler.DeleteSong))
		api.GET("/songs/:id/verses", adapter.ToGinHandler(handler.GetSongVerses))
	}

	return r
}
