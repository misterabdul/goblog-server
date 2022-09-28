package comments

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/service"
)

// @Tags        Comment (Writer)
// @Summary     Get My Comment
// @Description Get a comment of mine.
// @Router      /v1/auth/writer/comment/{uid} [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Comment's UID"
// @Success     200 {object} object{data=object{uid=string,postUid=string,parentCommentUid=string,email=string,name=string,content=string,replyCount=int,createdAt=string,deletedAt=string}}
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func GetMyComment(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			me              *models.UserModel
			comment         *models.CommentModel
			commentUid      primitive.ObjectID
			commentUidParam = c.Param("comment")
			err             error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if commentUid, err = primitive.ObjectIDFromHex(commentUidParam); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if comment, err = svc.Comment.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"postauthoruid": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": commentUid}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if comment == nil {
			responses.NotFound(c, errors.New("comment not found"))
			return
		}

		responses.AuthorizedComment(c, comment)
	}
}

// @Tags        Comment (Writer)
// @Summary     Get My Comments
// @Description Get comments of mine.
// @Router      /v1/auth/writer/comments [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Param       type  query    string false "Filter data by type, e.g.: ?type=trash, ?type=active."
// @Success     200   {object} object{data=[]object{uid=string,postUid=string,parentCommentUid=string,email=string,name=string,content=string,replyCount=int,createdAt=string,deletedAt=string}}
// @Success     204
// @Failure     401   {object} object{message=string}
// @Failure     404   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func GetMyComments(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			me          *models.UserModel
			comments    []*models.CommentModel
			queryParams = readCommonQueryParams(c)
			err         error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if comments, err = svc.Comment.GetMany(ctx, bson.M{
			"$and": append(queryParams,
				bson.M{"postauthoruid": bson.M{"$eq": me.UID}})},
			internalGin.GetFindOptions(c),
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(comments) == 0 {
			responses.NoContent(c)
			return
		}

		responses.AuthorizedComments(c, comments)
	}
}

// @Tags        Comment (Writer)
// @Summary     Get My Comments Stats
// @Description Get my comments's stats.
// @Router      /v1/auth/writer/comments/stats [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Param       type  query    string false "Filter data by type, e.g.: ?type=trash, ?type=active."
// @Success     200   {object} object{data=object{currentPage=int,totalPages=int,itemsPerPage=int,totalItems=int}}
// @Failure     401   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func GetMyCommentsStats(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			me          *models.UserModel
			count       int64
			queryParams = readCommonQueryParams(c)
			err         error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if count, err = svc.Comment.Count(ctx, bson.M{
			"$and": append(queryParams,
				bson.M{"postauthoruid": bson.M{"$eq": me.UID}})},
			internalGin.GetCountOptions(c),
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.ResourceStats(c, count)
	}
}

// @Tags        Comment (Writer)
// @Summary     Get My Post Comments
// @Description Get my post's comments.
// @Router      /v1/auth/writer/post/{uid}/comments [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid   path     string true "Post's UID or slug"
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Param       type  query    string false "Filter data by type, e.g.: ?type=trash, ?type=active."
// @Success     200   {object} object{data=[]object{uid=string,postUid=string,parentCommentUid=string,email=string,name=string,content=string,replyCount=int,createdAt=string,deletedAt=string}}
// @Success     204
// @Failure     401   {object} object{message=string}
// @Failure     404   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func GetMyPostComments(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			me           *models.UserModel
			comments     []*models.CommentModel
			postUid      primitive.ObjectID
			postUidParam = c.Param("post")
			queryParams  = readCommonQueryParams(c)
			err          error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if comments, err = svc.Comment.GetMany(ctx, bson.M{
			"$and": append(queryParams,
				bson.M{"postuid": bson.M{"$eq": postUid}},
				bson.M{"postauthoruid": bson.M{"$eq": me.UID}})},
			internalGin.GetFindOptions(c),
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(comments) == 0 {
			responses.NoContent(c)
			return
		}

		responses.AuthorizedComments(c, comments)
	}
}

// @Tags        Comment (Writer)
// @Summary     Get My Post Comments Stats
// @Description Get my post's comments stats.
// @Router      /v1/auth/writer/post/{uid}/comments/stats [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid   path     string true  "Post's UID or slug"
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Param       type  query    string false "Filter data by type, e.g.: ?type=trash, ?type=active."
// @Success     200   {object} object{data=object{currentPage=int,totalPages=int,itemsPerPage=int,totalItems=int}}
// @Failure     401   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func GetMyPostCommentsStats(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			me           *models.UserModel
			count        int64
			postUid      primitive.ObjectID
			postUidParam = c.Param("post")
			queryParams  = readCommonQueryParams(c)
			err          error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if count, err = svc.Comment.Count(ctx, bson.M{
			"$and": append(queryParams,
				bson.M{"postuid": bson.M{"$eq": postUid}},
				bson.M{"postauthoruid": bson.M{"$eq": me.UID}})},
			internalGin.GetCountOptions(c),
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.ResourceStats(c, count)
	}
}

// @Tags        Comment (Writer)
// @Summary     Delete My Comment (Soft)
// @Description Delete a comment of my post (soft-deleted).
// @Router      /v1/auth/writer/comment/{uid} [delete]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Comment's UID"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func TrashMyComment(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			me              *models.UserModel
			comment         *models.CommentModel
			parentComment   *models.CommentModel
			post            *models.PostModel
			commentUid      primitive.ObjectID
			commentUidParam = c.Param("comment")
			err             error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if commentUid, err = primitive.ObjectIDFromHex(commentUidParam); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if comment, err = svc.Comment.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"postauthoruid": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": commentUid}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if comment == nil {
			responses.NotFound(c, errors.New("comment not found"))
			return
		}
		if comment.ParentCommentUid == nil {
			if post, err = findCommentPost(c, svc, ctx, comment); err != nil {
				return
			}
			if err = svc.Comment.TrashOne(ctx, comment, post); err != nil {
				responses.InternalServerError(c, err)
				return
			}
		} else {
			if parentComment, err = findReplyParentComment(c, svc, ctx, comment); err != nil {
				return
			}
			if err = svc.Comment.TrashOneReply(ctx, comment, parentComment); err != nil {
				responses.InternalServerError(c, err)
				return
			}
		}

		responses.NoContent(c)
	}
}

// @Tags        Comment (Writer)
// @Summary     Restore My Comment (Soft)
// @Description Restore a deleted comment of my post (soft-deleted).
// @Router      /v1/auth/writer/comment/{uid}/detrash [put]
// @Router      /v1/auth/writer/comment/{uid}/detrash [patch]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Comment's UID"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DetrashMyComment(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			me              *models.UserModel
			comment         *models.CommentModel
			parentComment   *models.CommentModel
			post            *models.PostModel
			commentUid      primitive.ObjectID
			commentUidParam = c.Param("comment")
			err             error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if commentUid, err = primitive.ObjectIDFromHex(commentUidParam); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if comment, err = svc.Comment.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"postauthoruid": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": commentUid}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if comment == nil {
			responses.NotFound(c, errors.New("comment not found"))
			return
		}
		if comment.ParentCommentUid == nil {
			if post, err = findCommentPost(c, svc, ctx, comment); err != nil {
				return
			}
			if err = svc.Comment.RestoreOne(ctx, comment, post); err != nil {
				responses.InternalServerError(c, err)
				return
			}
		} else {
			if parentComment, err = findReplyParentComment(c, svc, ctx, comment); err != nil {
				return
			}
			if err = svc.Comment.RestoreOneReply(ctx, comment, parentComment); err != nil {
				responses.InternalServerError(c, err)
				return
			}
		}

		responses.NoContent(c)
	}
}

// @Tags        Comment (Writer)
// @Summary     Delete My Comment (Permanent)
// @Description Delete a comment of my post (permanent).
// @Router      /v1/auth/writer/comment/{uid}/permanent [delete]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Comment's UID"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DeleteMyComment(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			me              *models.UserModel
			comment         *models.CommentModel
			parentComment   *models.CommentModel
			post            *models.PostModel
			commentUid      primitive.ObjectID
			commentUidParam = c.Param("comment")
			err             error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if commentUid, err = primitive.ObjectIDFromHex(commentUidParam); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if comment, err = svc.Comment.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"postauthoruid": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": commentUid}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if comment == nil {
			responses.NotFound(c, errors.New("comment not found"))
			return
		}
		if comment.ParentCommentUid == nil {
			if post, err = findCommentPost(c, svc, ctx, comment); err != nil {
				return
			}
			if err = svc.Comment.DeleteOne(ctx, comment, post); err != nil {
				responses.InternalServerError(c, err)
				return
			}
		} else {
			if parentComment, err = findReplyParentComment(c, svc, ctx, comment); err != nil {
				return
			}
			if err = svc.Comment.DeleteOneReply(ctx, comment, parentComment); err != nil {
				responses.InternalServerError(c, err)
				return
			}
		}

		responses.NoContent(c)
	}
}

func findCommentPost(
	c *gin.Context,
	svc *service.Service,
	ctx context.Context,
	comment *models.CommentModel,
) (post *models.PostModel, err error) {

	if post, err = svc.Post.GetOne(ctx, bson.M{
		"$and": []bson.M{
			{"deletedat": bson.M{"$eq": primitive.Null{}}},
			{"_id": bson.M{"$eq": comment.PostUid}}}},
	); err != nil {
		responses.InternalServerError(c, err)
		return nil, err
	}
	if post == nil {
		responses.NotFound(c, errors.New("post not found"))
		return nil, err
	}

	return post, nil
}

func findReplyParentComment(
	c *gin.Context,
	svc *service.Service,
	ctx context.Context,
	reply *models.CommentModel,
) (comment *models.CommentModel, err error) {
	if comment, err = svc.Comment.GetOne(ctx, bson.M{
		"$and": []bson.M{
			{"deletedat": bson.M{"$eq": primitive.Null{}}},
			{"_id": bson.M{"$eq": reply.ParentCommentUid}}}},
	); err != nil {
		responses.InternalServerError(c, err)
		return nil, err
	}
	if comment == nil {
		responses.NotFound(c, errors.New("parent comment not found"))
		return nil, err
	}

	return comment, nil
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
	default:
		extraQuery = append(extraQuery,
			bson.M{"deletedat": bson.M{"$eq": primitive.Null{}}})
	}

	return extraQuery
}
