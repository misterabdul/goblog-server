package users

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/handlers/helpers"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

// Get single user record publicly
func GetPublicUser(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			user      *models.UserModel
			userId    primitive.ObjectID
			userQuery = c.Param("user")
			err       error
		)

		if userId, err = primitive.ObjectIDFromHex(userQuery); err != nil {
			userId = primitive.ObjectID{}
		}
		if user, err = repositories.GetUser(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"$or": []bson.M{
					{"_id": userId},
					{"username": userQuery}}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if user == nil {
			responses.NotFound(c, errors.New("user not found"))
			return
		}

		responses.PublicUser(c, user)
	}
}

// Get multiple user records publicly
func GetPublicUsers(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			users []*models.UserModel
			err   error
		)

		if users, err = repositories.GetUsers(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}}},
		}, helpers.GetFindOptions(c)); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		log.Println(users)
		if len(users) == 0 {
			responses.NoContent(c)
			return
		}

		responses.PublicUsers(c, users)
	}
}
