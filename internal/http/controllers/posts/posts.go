package posts

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetPublicPost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{})
	}
}

func GetPublicPostSlug(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{})
	}
}

func GetPublicPosts(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{})
	}
}
