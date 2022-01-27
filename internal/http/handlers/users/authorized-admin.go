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
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

func GetUser(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			userService  = service.New(c, ctx, dbConn)
			user         *models.UserModel
			userUid      primitive.ObjectID
			userUidParam = c.Param("user")
			err          error
		)

		defer cancel()
		if userUid, err = primitive.ObjectIDFromHex(userUidParam); err != nil {
			responses.IncorrectUserId(c, err)
			return
		}
		if user, err = userService.GetUser(bson.M{
			"_id": bson.M{"$eq": userUid},
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
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			userService = service.New(c, ctx, dbConn)
			users       []*models.UserModel
			typeParam   = c.DefaultQuery("type", "active")
			extraQuery  = []bson.M{}
			err         error
		)

		defer cancel()
		switch true {
		case typeParam == "trash":
			extraQuery = append(extraQuery,
				bson.M{"deletedat": bson.M{"$ne": primitive.Null{}}})
		case typeParam == "active":
			fallthrough
		default:
			extraQuery = append(extraQuery,
				bson.M{"deletedat": bson.M{"$eq": primitive.Null{}}})
		}
		if users, err = userService.GetUsers(bson.M{
			"$and": extraQuery,
		}); err != nil {
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
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			userService = service.New(c, ctx, dbConn)
			me          *models.UserModel
			newUser     *models.UserModel
			form        *forms.CreateUserForm
			err         error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if form, err = requests.GetCreateUserForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if err = form.Validate(userService, me); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if newUser, err = form.ToUserModel(); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = userService.CreateUser(newUser); err != nil {
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
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			userService  = service.New(c, ctx, dbConn)
			me           *models.UserModel
			user         *models.UserModel
			userUid      primitive.ObjectID
			userUidParam = c.Param("user")
			form         *forms.UpdateUserForm
			err          error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if userUid, err = primitive.ObjectIDFromHex(userUidParam); err != nil {
			responses.IncorrectUserId(c, err)
			return
		}
		if user, err = userService.GetUser(bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": bson.M{"$eq": userUid}}},
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
		if err = form.Validate(userService, me, user); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if user, err = form.ToUserModel(user); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = userService.UpdateUser(user); err != nil {
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
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			userService  = service.New(c, ctx, dbConn)
			user         *models.UserModel
			userUid      primitive.ObjectID
			userUidParam = c.Param("user")
			err          error
		)

		defer cancel()
		if userUid, err = primitive.ObjectIDFromHex(userUidParam); err != nil {
			responses.IncorrectUserId(c, err)
			return
		}
		if user, err = userService.GetUser(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": userUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if user == nil {
			responses.NotFound(c, errors.New("user not found"))
			return
		}
		if err = userService.TrashUser(user); err != nil {
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
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			userService  = service.New(c, ctx, dbConn)
			user         *models.UserModel
			userUid      primitive.ObjectID
			userUidParam = c.Param("user")
			err          error
		)

		defer cancel()
		if userUid, err = primitive.ObjectIDFromHex(userUidParam); err != nil {
			responses.IncorrectUserId(c, err)
			return
		}
		if user, err = userService.GetUser(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": bson.M{"$eq": userUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if user == nil {
			responses.NotFound(c, errors.New("user not found"))
			return
		}
		if err = userService.DetrashUser(user); err != nil {
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
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			userService  = service.New(c, ctx, dbConn)
			user         *models.UserModel
			userUid      primitive.ObjectID
			userUidParam = c.Param("user")
			err          error
		)

		defer cancel()
		if userUid, err = primitive.ObjectIDFromHex(userUidParam); err != nil {
			responses.IncorrectUserId(c, err)
			return
		}
		if user, err = userService.GetUser(bson.M{
			"_id": bson.M{"$eq": userUid},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if user == nil {
			responses.NotFound(c, errors.New("user not found"))
			return
		}
		if err = userService.DeleteUser(user); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}
