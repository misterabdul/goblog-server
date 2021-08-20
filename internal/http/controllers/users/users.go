package users

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

// Get single user record publicly
func GetPublicUser(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var dbConn *mongo.Database
		var userData *models.UserModel
		var userId primitive.ObjectID
		var err error
		userIdQuery := c.Param("user")

		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if userId, err = primitive.ObjectIDFromHex(userIdQuery); err != nil {
			responses.NotFound(c, err)
			return
		}
		if userData, err = repositories.GetUser(ctx, dbConn, bson.M{"_id": userId}); err != nil {
			responses.NotFound(c, err)
			return
		}

		responses.PublicUser(c, userData)
	}
}

// Get multiple user records publicly
func GetPublicUsers(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var dbConn *mongo.Database
		var usersData []*models.UserModel
		var err error

		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if usersData, err = repositories.GetUsers(ctx, dbConn, bson.M{}, 10, "createdAt", false); err != nil {
			responses.NotFound(c, err)
			return
		}

		responses.PublicUsers(c, usersData)
	}
}
