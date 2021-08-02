package authorize

import "github.com/gin-gonic/gin"

func Authorize(level string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
