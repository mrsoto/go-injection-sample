package router

import (
	"github.com/gin-gonic/gin"
)

type GetAbums = gin.HandlerFunc
type GetAlbumByID = gin.HandlerFunc
type PostAlbums = gin.HandlerFunc

type Controller struct {
	GetAbums
	GetAlbumByID
	PostAlbums
}

func RegisterRoutes(router *gin.Engine, c Controller) {
	router.GET("/albums", c.GetAbums)
	router.GET("/albums/:id", c.GetAlbumByID)
	router.POST("/albums", c.PostAlbums)
}
