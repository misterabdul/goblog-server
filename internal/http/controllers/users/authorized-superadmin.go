package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminizeUser(c *gin.Context) {

	c.JSON(http.StatusNotImplemented, gin.H{})
}

func DeadminizeUser(c *gin.Context) {

	c.JSON(http.StatusNotImplemented, gin.H{})
}
