package main

import (
	"log"
	"mediaSequencer/backend/internal/handler"
	"mediaSequencer/backend/internal/repository"
	"mediaSequencer/backend/internal/service"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Overload(".env")
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db := repository.ConnectDatabase()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	windowRepo := repository.NewWindowRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	playlistRepo := repository.NewPlaylistRepository(db)

	windowHandler := handler.NewWindowHandler(windowRepo)
	mediaHandler := handler.NewMediaHandler(mediaRepo)
	playlistHandler := handler.NewPlaylistHandler(playlistRepo)

	router.GET("/windows", windowHandler.GetWindows)

	router.GET("/media", mediaHandler.GetMedia)
	router.POST("/media", mediaHandler.CreateMedia)

	router.GET("/windows/:windowId/playlist", playlistHandler.GetPlaylist)
	router.POST("/playlist", playlistHandler.AddToPlaylist)

	syncService := service.NewSyncService()
	syncHandler := handler.NewSyncHandler(syncService)

	router.POST("/sync", syncHandler.StartSync)
	router.GET("/sync", syncHandler.GetSync)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}
