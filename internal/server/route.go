package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Initialize all routes.
func initRoute(server *gin.Engine) {
	api := server.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})
		}
	}
	server.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
	})
}
