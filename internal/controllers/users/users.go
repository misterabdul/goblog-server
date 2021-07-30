package users

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/internal/responses"
)

// Get single user record publicly
func GetPublicUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var dbConn *mongo.Database
	var userData *models.UserModel
	var userId primitive.ObjectID
	var err error
	userIdQuery := c.Param("user")

	if dbConn, err = repositories.GetDBConnDefault(ctx); err != nil {
		responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if userId, err = primitive.ObjectIDFromHex(userIdQuery); err != nil {
		responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if userData, err = repositories.GetUser(ctx, dbConn, bson.M{"_id": userId}); err != nil {
		responses.Basic(c, http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	responses.User(c, userData)
}

// Get multiple user records publicly
func GetPublicUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var dbConn *mongo.Database
	var usersData []*models.UserModel
	var err error

	if dbConn, err = repositories.GetDBConnDefault(ctx); err != nil {
		responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if usersData, err = repositories.GetUsers(ctx, dbConn, bson.M{}, 10, "createdAt", false); err != nil {
		responses.Basic(c, http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	responses.Users(c, usersData)
}
