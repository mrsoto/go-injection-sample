package main

import (
	"example/web-service-gin/album"
	"example/web-service-gin/album/private/persistence"
	"example/web-service-gin/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.LstdFlags)
	r := gin.Default()

	repository := persistence.NewInMemoryRepository()
	albumCtrl := album.NewController(repository)
	router.RegisterRoutes(r, albumCtrl)

	if err := r.Run("localhost:8080"); err != nil {
		log.Panic("unable to start the server")
	}
}
