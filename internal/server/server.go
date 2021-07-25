package server

import "github.com/gin-gonic/gin"

// Get the server engine.
func GetServer() *gin.Engine {
	ginEngine := gin.Default()

	initRoute(ginEngine)

	return ginEngine
}
