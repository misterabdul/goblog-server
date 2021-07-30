package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {

	c.JSON(http.StatusNotImplemented, gin.H{})
}

func UpdateMe(c *gin.Context) {

	c.JSON(http.StatusNotImplemented, gin.H{})
}

func UpdateMePassword(c *gin.Context) {

	c.JSON(http.StatusNotImplemented, gin.H{})
}
