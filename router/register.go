package router

import (
	"github.com/gin-gonic/gin"
)

type Controller struct {
	GetAbums     gin.HandlerFunc
	GetAlbumByID gin.HandlerFunc
	PostAlbums   gin.HandlerFunc
}

func RegisterRoutes(router *gin.Engine, c Controller) {
	router.GET("/albums", c.GetAbums)
	router.GET("/albums/:id", c.GetAlbumByID)
	router.POST("/albums", c.PostAlbums)
}
