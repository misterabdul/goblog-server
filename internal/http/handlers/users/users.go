package users

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

// Get single user record publicly
func GetPublicUser(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			userService = service.NewUserService(c, ctx, dbConn)
			user        *models.UserModel
			userUid     interface{}
			userParam   = c.Param("user")
			err         error
		)

		defer cancel()
		if userUid, err = primitive.ObjectIDFromHex(userParam); err != nil {
			userUid = nil
		}
		if user, err = userService.GetUser(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"$or": []bson.M{
					{"_id": bson.M{"$eq": userUid}},
					{"username": bson.M{"$eq": userParam}}}}},
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
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			userService = service.NewUserService(c, ctx, dbConn)
			users       []*models.UserModel
			err         error
		)

		defer cancel()
		if users, err = userService.GetUsers(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(users) == 0 {
			responses.NoContent(c)
			return
		}

		responses.PublicUsers(c, users)
	}
}
