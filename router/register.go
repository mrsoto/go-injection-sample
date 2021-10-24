package router

import (
	"example/web-service-gin/album"

	"github.com/gin-gonic/gin"
)

type Services struct {
	GetAbums     gin.HandlerFunc
	GetAlbumByID gin.HandlerFunc
	PostAlbums   gin.HandlerFunc
}

func RegisterRoutes(router *gin.Engine, s Services) {
	router.GET("/albums", s.GetAbums)
	router.GET("/albums/:id", s.GetAlbumByID)
	router.POST("/albums", album.PostAlbums)
}
