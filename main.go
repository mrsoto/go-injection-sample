package main

import (
	"example/web-service-gin/album"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.LstdFlags)
	router := gin.Default()
	router.GET("/albums", album.GetAlbums)
	router.GET("/albums/:id", album.GetAlbumByID)
	router.POST("/albums", album.PostAlbums)

	if err := router.Run("localhost:8080"); err != nil {
		log.Panic("unable to start the server")
	}
}
