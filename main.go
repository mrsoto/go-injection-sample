package main

import (
	"example/web-service-gin/album"
	repository "example/web-service-gin/album/private/persistence"
	"example/web-service-gin/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.LstdFlags)
	r := gin.Default()

	repository := repository.NewInMemoryRepository()
	albumSrvs := album.NewServices(repository)

	router.RegisterRoutes(r, router.Controller{
		GetAbums:     albumSrvs.GetAlbums,
		GetAlbumByID: albumSrvs.GetAlbumByID,
		PostAlbums:   albumSrvs.PostAlbums,
	})

	if err := r.Run("localhost:8080"); err != nil {
		log.Panic("unable to start the server")
	}
}
