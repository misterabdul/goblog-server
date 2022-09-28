package comments

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

// @Tags        Comment (Editor)
// @Summary     Get Comment
// @Description Get a comment.
// @Router      /v1/auth/editor/comment/{uid} [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Comment's UID"
// @Success     200 {object} object{data=object{uid=string,postUid=string,parentCommentUid=string,email=string,name=string,content=string,replyCount=int,createdAt=string,deletedAt=string}}
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func GetComment(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			comment         *models.CommentModel
			commentUid      primitive.ObjectID
			commentUidParam = c.Param("comment")
			err             error
		)

		defer cancel()
		if commentUid, err = primitive.ObjectIDFromHex(commentUidParam); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if comment, err = svc.Comment.GetOne(ctx, bson.M{
			"_id": bson.M{"$eq": commentUid}},
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

// @Tags        Comment (Editor)
// @Summary     Get Comments
// @Description Get comments.
// @Router      /v1/auth/editor/comments [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Param       type  query    string false "Filter data by type, e.g.: ?type=trash, ?type=published, ?type=draft."
// @Success     200   {object} object{data=[]object{uid=string,postUid=string,parentCommentUid=string,email=string,name=string,content=string,replyCount=int,createdAt=string,deletedAt=string}}
// @Success     204
// @Failure     401   {object} object{message=string}
// @Failure     404   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func GetComments(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			comments    []*models.CommentModel
			queryParams = readCommonQueryParams(c)
			err         error
		)

		defer cancel()
		if comments, err = svc.Comment.GetMany(ctx, bson.M{
			"$and": queryParams},
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

// @Tags        Comment (Editor)
// @Summary     Get Comments Stats
// @Description Get comments's stats.
// @Router      /v1/auth/editor/comments/stats [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Param       type  query    string false "Filter data by type, e.g.: ?type=trash, ?type=published, ?type=draft."
// @Success     200   {object} object{data=object{currentPage=int,totalPages=int,itemsPerPage=int,totalItems=int}}
// @Failure     401   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func GetCommentsStats(
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
		if count, err = svc.Comment.Count(ctx, bson.M{
			"$and": queryParams},
			internalGin.GetCountOptions(c),
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.ResourceStats(c, count)
	}
}

// @Tags        Comment (Editor)
// @Summary     Get Post Comments
// @Description Get post's comments.
// @Router      /v1/auth/editor/post/{uid}/comments [get]
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
func GetPostComments(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			comments     []*models.CommentModel
			postUid      primitive.ObjectID
			postUidParam = c.Param("post")
			queryParams  = readCommonQueryParams(c)
			err          error
		)

		defer cancel()
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if comments, err = svc.Comment.GetMany(ctx, bson.M{
			"$and": append(queryParams,
				bson.M{"postuid": bson.M{"$eq": postUid}})},
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

// @Tags        Comment (Editor)
// @Summary     Get Post Comments Stats
// @Description Get post's comments stats.
// @Router      /v1/auth/editor/post/{uid}/comments/stats [get]
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
func GetPostCommentsStats(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			count        int64
			postUid      primitive.ObjectID
			postUidParam = c.Param("post")
			queryParams  = readCommonQueryParams(c)
			err          error
		)

		defer cancel()
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if count, err = svc.Comment.Count(ctx, bson.M{
			"$and": append(queryParams,
				bson.M{"postuid": bson.M{"$eq": postUid}})},
			internalGin.GetCountOptions(c),
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.ResourceStats(c, count)
	}
}

// @Tags        Comment (Editor)
// @Summary     Delete Comment (Soft)
// @Description Delete a comment (soft-deleted).
// @Router      /v1/auth/editor/comment/{uid} [delete]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Comment's UID"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func TrashComment(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel      = context.WithTimeout(context.Background(), maxCtxDuration)
			comment          *models.CommentModel
			parentComment    *models.CommentModel
			post             *models.PostModel
			commentUid       primitive.ObjectID
			commentdUidParam = c.Param("comment")
			err              error
		)

		defer cancel()
		if commentUid, err = primitive.ObjectIDFromHex(commentdUidParam); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if comment, err = svc.Comment.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
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

// @Tags        Comment (Editor)
// @Summary     Restore Comment (Soft)
// @Description Restore deleted comment (soft-deleted).
// @Router      /v1/auth/editor/comment/{uid}/detrash [put]
// @Router      /v1/auth/editor/comment/{uid}/detrash [patch]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Comment's UID"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DetrashComment(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			comment         *models.CommentModel
			parentComment   *models.CommentModel
			post            *models.PostModel
			commentUid      primitive.ObjectID
			commentUidParam = c.Param("comment")
			err             error
		)

		defer cancel()
		if commentUid, err = primitive.ObjectIDFromHex(commentUidParam); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if comment, err = svc.Comment.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
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

// @Tags        Comment (Editor)
// @Summary     Delete Comment (Permanent)
// @Description Delete comment (permanent).
// @Router      /v1/auth/editor/comment/{uid}/permanent [delete]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Comment's UID"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DeleteComment(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			comment         *models.CommentModel
			parentComment   *models.CommentModel
			post            *models.PostModel
			commentUid      primitive.ObjectID
			commentUidParam = c.Param("comment")
			err             error
		)

		defer cancel()
		if commentUid, err = primitive.ObjectIDFromHex(commentUidParam); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if comment, err = svc.Comment.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"_id": bson.M{"$eq": commentUid}}},
		}); err != nil {
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
