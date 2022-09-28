package users

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/service"
)

// @Tags        User (Public)
// @Summary     Get Public User
// @Description Get a user that available publicly.
// @Router      /v1/user/{uid} [get]
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "User's UID or slug"
// @Success     200 {object} object{data=object{uid=string,username=string,email=string,firstName=string,lastName=string}}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func GetPublicUser(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			user        *models.UserModel
			userUid     interface{}
			userParam   = c.Param("user")
			err         error
		)

		defer cancel()
		if userUid, err = primitive.ObjectIDFromHex(userParam); err != nil {
			userUid = nil
		}
		if user, err = svc.User.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"$or": []bson.M{
					{"_id": bson.M{"$eq": userUid}},
					{"username": bson.M{"$eq": userParam}}}}}},
		); err != nil {
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

// @Tags        Users (Public)
// @Summary     Get Public Users
// @Description Get users that available publicly.
// @Router      /v1/users [get]
// @Produce     application/json
// @Produce     application/msgpack
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Success     200   {object} object{data=[]object{uid=string,username=string,email=string,firstName=string,lastName=string}}
// @Success     204
// @Failure     500   {object} object{message=string}
func GetPublicUsers(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			users       []*models.UserModel
			err         error
		)

		defer cancel()
		if users, err = svc.User.GetMany(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}}}},
			internalGin.GetFindOptions(c),
		); err != nil {
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
