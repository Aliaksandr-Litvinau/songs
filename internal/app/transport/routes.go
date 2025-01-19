package transport

import (
	_ "songs/docs"
	"songs/internal/app/middleware"
	"songs/internal/app/transport/adapter"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(svc SongService) *gin.Engine {
	r := gin.Default()

	// Add Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// Add Prometheus middleware
	r.Use(middleware.PrometheusMiddleware())

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
