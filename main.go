package main

import (
	"example/web-service-gin/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.LstdFlags)
	r := gin.Default()
	router.RegisterRoutes(r)

	if err := r.Run("localhost:8080"); err != nil {
		log.Panic("unable to start the server")
	}
}
