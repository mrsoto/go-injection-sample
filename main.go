package main

import (
	"example/web-service-gin/album"
	"example/web-service-gin/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.LstdFlags)
	r := gin.Default()
	albumSrvs := router.Services{
		GetAbums:     album.GetAlbums,
		GetAlbumByID: album.GetAlbumByID,
		PostAlbums:   album.PostAlbums,
	}
	router.RegisterRoutes(r, albumSrvs)

	if err := r.Run("localhost:8080"); err != nil {
		log.Panic("unable to start the server")
	}
}
