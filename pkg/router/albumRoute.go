package router

import (
	"github.com/gin-gonic/gin"
)

type Albums interface {
	GetAlbums(*gin.Context)
	GetAlbumByID(*gin.Context)
	PostAlbums(*gin.Context)
}

func RegisterRoutes(router *gin.Engine, c Albums) {
	router.GET("/albums", c.GetAlbums)
	router.GET("/albums/:id", c.GetAlbumByID)
	router.POST("/albums", c.PostAlbums)
}
