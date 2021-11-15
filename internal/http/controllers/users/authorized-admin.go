package users

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUsers(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{})
	}
}

func GetUser(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{})
	}
}

func CreateUser(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{})
	}
}

func UpdateUser(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{})
	}
}

func TrashUser(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{})
	}
}

func DeleteUser(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{})
	}
}
