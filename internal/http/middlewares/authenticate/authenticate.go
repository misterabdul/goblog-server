package authenticate

import "github.com/gin-gonic/gin"

// Check the authentication status of given user.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
