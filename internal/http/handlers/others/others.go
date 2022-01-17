package others

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/responses"
)

func Ping() gin.HandlerFunc {

	return func(c *gin.Context) {
		responses.Success(c, gin.H{
			"message": "pong"})
	}
}

func NotFound() gin.HandlerFunc {

	return func(c *gin.Context) {
		responses.NotFound(c, nil)
	}
}
