package router

import (
	"example/web-service-gin/album"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/albums", album.GetAlbums)
	router.GET("/albums/:id", album.GetAlbumByID)
	router.POST("/albums", album.PostAlbums)
}
