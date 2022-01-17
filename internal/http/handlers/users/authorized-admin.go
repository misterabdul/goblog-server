package users

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/handlers/helpers"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

func GetUser(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			user        *models.UserModel
			userId      primitive.ObjectID
			userIdQuery = c.Param("user")
			err         error
		)

		if userId, err = primitive.ObjectIDFromHex(userIdQuery); err != nil {
			responses.IncorrectUserId(c, err)
			return
		}
		if user, err = repositories.GetUser(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"_id": userId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if user == nil {
			responses.NotFound(c, errors.New("user not found"))
			return
		}

		responses.AuthorizedUser(c, user)
	}
}

func GetUsers(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			users     []*models.UserModel
			typeParam = c.DefaultQuery("type", "draft")
			typeQuery []bson.M
			err       error
		)

		switch true {
		case typeParam == "trash":
			typeQuery = []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}}}
		case typeParam == "active":
			fallthrough
		default:
			typeQuery = []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}}}
		}

		if users, err = repositories.GetUsers(ctx, dbConn, bson.M{
			"$and": typeQuery,
		}, helpers.GetFindOptions(c)); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(users) == 0 {
			responses.NoContent(c)
			return
		}

		responses.AuthorizedUsers(c, users)
	}
}

func CreateUser(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me       *models.UserModel
			newUser  *models.UserModel
			form     *forms.CreateUserForm
			err      error
			writeErr mongo.WriteException
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if form, err = requests.GetCreateUserForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if err = form.Validate(ctx, dbConn, me); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if newUser, err = form.ToUserModel(); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = repositories.CreateUser(ctx, dbConn, newUser); err != nil {
			if errors.As(err, &writeErr) {
				responses.FormIncorrect(c, err)
				return
			}
			responses.InternalServerError(c, err)
			return
		}

		responses.AuthorizedUser(c, newUser)
	}
}

func UpdateUser(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me          *models.UserModel
			user        *models.UserModel
			userId      primitive.ObjectID
			userIdQuery = c.Param("user")
			form        *forms.UpdateUserForm
			err         error
			writeErr    mongo.WriteException
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if userId, err = primitive.ObjectIDFromHex(userIdQuery); err != nil {
			responses.IncorrectUserId(c, err)
			return
		}
		if user, err = repositories.GetUser(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": userId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if user == nil {
			responses.NotFound(c, errors.New("user not found"))
			return
		}
		if form, err = requests.GetUpdateUserForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if err = form.Validate(ctx, dbConn, me, user); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if user, err = form.ToUserModel(user); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = repositories.UpdateUser(ctx, dbConn, user); err != nil {
			if errors.As(err, &writeErr) {
				responses.FormIncorrect(c, err)
				return
			}
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func TrashUser(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			user        *models.UserModel
			userId      primitive.ObjectID
			userIdQuery = c.Param("user")
			err         error
		)

		if userId, err = primitive.ObjectIDFromHex(userIdQuery); err != nil {
			responses.IncorrectUserId(c, err)
			return
		}
		if user, err = repositories.GetUser(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": userId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if user == nil {
			responses.NotFound(c, errors.New("user not found"))
			return
		}
		if err = repositories.TrashUser(ctx, dbConn, user); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DetrashUser(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			user        *models.UserModel
			userId      primitive.ObjectID
			userIdQuery = c.Param("user")
			err         error
		)

		if userId, err = primitive.ObjectIDFromHex(userIdQuery); err != nil {
			responses.IncorrectUserId(c, err)
			return
		}
		if user, err = repositories.GetUser(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": userId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if user == nil {
			responses.NotFound(c, errors.New("user not found"))
			return
		}
		if err = repositories.DetrashUser(ctx, dbConn, user); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DeleteUser(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			user        *models.UserModel
			userId      primitive.ObjectID
			userIdQuery = c.Param("user")
			err         error
		)

		if userId, err = primitive.ObjectIDFromHex(userIdQuery); err != nil {
			responses.IncorrectUserId(c, err)
			return
		}
		if user, err = repositories.GetUser(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"_id": userId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if user == nil {
			responses.NotFound(c, errors.New("user not found"))
			return
		}
		if err = repositories.DeleteUser(ctx, dbConn, user); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}
