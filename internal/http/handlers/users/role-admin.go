package users

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/service"
)

// @Tags        User (Admin)
// @Summary     Get User
// @Description Get a user.
// @Router      /v1/auth/admin/user/{uid} [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "User's UID"
// @Success     200 {object} object{data=object{uid=string,username=string,email=string,firstName=string,lastName=string,roles=[]object{level=int,name=string,since=string},createdAt=time,updatedAt=time}}
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func GetUser(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
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
		if user, err = svc.User.GetOne(ctx, bson.M{
			"_id": bson.M{"$eq": userUid}},
		); err != nil {
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

// @Tags        User (Admin)
// @Summary     Get Users
// @Description Get users.
// @Router      /v1/auth/admin/users [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Success     200   {object} object{data=[]object{uid=string,username=string,email=string,firstName=string,lastName=string,roles=[]object{level=int,name=string,since=string},createdAt=time,updatedAt=time}}
// @Success     204
// @Failure     401   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func GetUsers(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			users       []*models.UserModel
			queryParams = readCommonQueryParams(c)
			err         error
		)

		defer cancel()
		if users, err = svc.User.GetMany(ctx, bson.M{
			"$and": queryParams},
			internalGin.GetFindOptions(c),
		); err != nil {
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

// @Tags        User (Admin)
// @Summary     Get Users Stats
// @Description Get users's stats.
// @Router      /v1/auth/admin/users/stats [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Success     200   {object} object{data=object{currentPage=int,totalPages=int,itemsPerPage=int,totalItems=int}}
// @Failure     401   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func GetUsersStats(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			count       int64
			queryParams = readCommonQueryParams(c)
			err         error
		)

		defer cancel()
		if count, err = svc.User.Count(ctx, bson.M{
			"$and": queryParams},
			internalGin.GetCountOptions(c),
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.ResourceStats(c, count)
	}
}

// @Tags        User (Admin)
// @Summary     Create User
// @Description Create a new user.
// @Router      /v1/auth/admin/user [post]
// @Security    BearerAuth
// @Accept      application/json
// @Accept      application/msgpack
// @Produce     application/json
// @Produce     application/msgpack
// @Param       form body     object{firstName=string,lastName=string,username=string,email=string,password=string,passwordConfirm=string,roles=[]int} true "Create user form"
// @Success     200  {object} object{data=object{uid=string,username=string,email=string,firstName=string,lastName=string,roles=[]object{level=int,name=string,since=string},createdAt=time,updatedAt=time}}
// @Failure     401  {object} object{message=string}
// @Failure     422  {object} object{message=string}
// @Failure     500  {object} object{message=string}
func CreateUser(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
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
		if err = form.Validate(svc, ctx, me); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if newUser, err = form.ToUserModel(); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = svc.User.SaveOne(ctx, newUser); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.AuthorizedUser(c, newUser)
	}
}

// @Tags        User (Admin)
// @Summary     Update User
// @Description Update User
// @Router      /v1/auth/admin/user/{uid} [put]
// @Router      /v1/auth/admin/user/{uid} [patch]
// @Security    BearerAuth
// @Accept      application/json
// @Accept      application/msgpack
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid  path     string                                                                                                                   true "User's UID"
// @Param       form body     object{firstName=string,lastName=string,username=string,email=string,password=string,passwordConfirm=string,roles=[]int} true "Update user form"
// @Success     204
// @Failure     401  {object} object{message=string}
// @Failure     404  {object} object{message=string}
// @Failure     422  {object} object{message=string}
// @Failure     500  {object} object{message=string}
func UpdateUser(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
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
		if user, err = svc.User.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": bson.M{"$eq": userUid}}}},
		); err != nil {
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
		if err = form.Validate(svc, ctx, me, user); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if user, err = form.ToUserModel(user); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = svc.User.UpdateOne(ctx, user); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        User (Admin)
// @Summary     Delete User (Soft)
// @Description Delete a user (soft-deleted).
// @Router      /v1/auth/admin/user/{uid} [delete]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "User's UID"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     422 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func TrashUser(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
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
		if user, err = svc.User.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": userUid}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if user == nil {
			responses.NotFound(c, errors.New("user not found"))
			return
		}
		if err = svc.User.TrashOne(ctx, user); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        User (Admin)
// @Summary     Restore user (Soft)
// @Description Restore a deleted user (soft-deleted).
// @Router      /v1/auth/admin/user/{uid}/detrash [put]
// @Router      /v1/auth/admin/user/{uid}/detrash [patch]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "User's UID"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     422 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DetrashUser(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
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
		if user, err = svc.User.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": bson.M{"$eq": userUid}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if user == nil {
			responses.NotFound(c, errors.New("user not found"))
			return
		}
		if err = svc.User.RestoreOne(ctx, user); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        User (Admin)
// @Summary     Delete User (Permanent)
// @Description Delete a user (permanent).
// @Router      /v1/auth/admin/user/{uid}/permanent [delete]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "User's UID"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     422 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DeleteUser(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
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
		if user, err = svc.User.GetOne(ctx, bson.M{
			"_id": bson.M{"$eq": userUid}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if user == nil {
			responses.NotFound(c, errors.New("user not found"))
			return
		}
		if err = svc.User.DeleteOne(ctx, user); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func readCommonQueryParams(c *gin.Context) []bson.M {
	var (
		typeParam  = c.DefaultQuery("type", "active")
		extraQuery = []bson.M{}
	)

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

	return extraQuery
}
