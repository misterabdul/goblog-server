package comments

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

func GetMyComment(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService  = service.New(c, ctx, dbConn)
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
		if comment, err = commentService.GetComment(bson.M{
			"$and": []bson.M{
				{"postauthoruid": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": commentUid}}},
		}); err != nil {
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

func GetMyComments(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel    = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService = service.New(c, ctx, dbConn)
			me             *models.UserModel
			comments       []*models.CommentModel
			queryParams    = readCommonQueryParams(c)
			err            error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if comments, err = commentService.GetComments(bson.M{
			"$and": append(queryParams,
				bson.M{"postauthoruid": bson.M{"$eq": me.UID}}),
		}, false); err != nil {
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

func GetMyCommentsStats(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel    = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService = service.New(c, ctx, dbConn)
			me             *models.UserModel
			count          int64
			queryParams    = readCommonQueryParams(c)
			err            error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if count, err = commentService.GetCommentCount(bson.M{
			"$and": append(queryParams,
				bson.M{"postauthoruid": bson.M{"$eq": me.UID}}),
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.ResourceStats(c, count)
	}
}

func GetMyPostComments(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel    = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService = service.New(c, ctx, dbConn)
			me             *models.UserModel
			comments       []*models.CommentModel
			postUid        primitive.ObjectID
			postUidParam   = c.Param("post")
			queryParams    = readCommonQueryParams(c)
			err            error
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
		if comments, err = commentService.GetComments(bson.M{
			"$and": append(queryParams,
				bson.M{"postuid": bson.M{"$eq": postUid}},
				bson.M{"postauthoruid": bson.M{"$eq": me.UID}}),
		}, false); err != nil {
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

func GetMyPostCommentsStats(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel    = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService = service.New(c, ctx, dbConn)
			me             *models.UserModel
			count          int64
			postUid        primitive.ObjectID
			postUidParam   = c.Param("post")
			queryParams    = readCommonQueryParams(c)
			err            error
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
		if count, err = commentService.GetCommentCount(bson.M{
			"$and": append(queryParams,
				bson.M{"postuid": bson.M{"$eq": postUid}},
				bson.M{"postauthoruid": bson.M{"$eq": me.UID}}),
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.ResourceStats(c, count)
	}
}

func TrashMyComment(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService  = service.New(c, ctx, dbConn)
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
		if comment, err = commentService.GetComment(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"postauthoruid": bson.M{"$eq": me.UID}},
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
			if post, err = findCommentPost(c, commentService, comment); err != nil {
				return
			}
			if err = commentService.TrashComment(comment, post); err != nil {
				responses.InternalServerError(c, err)
				return
			}
		} else {
			if parentComment, err = findReplyParentComment(c, commentService, comment); err != nil {
				return
			}
			if err = commentService.TrashCommentReply(comment, parentComment); err != nil {
				responses.InternalServerError(c, err)
				return
			}
		}

		responses.NoContent(c)
	}
}

func DetrashMyComment(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService  = service.New(c, ctx, dbConn)
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
		if comment, err = commentService.GetComment(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"postauthoruid": bson.M{"$eq": me.UID}},
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
			if post, err = findCommentPost(c, commentService, comment); err != nil {
				return
			}
			if err = commentService.DetrashComment(comment, post); err != nil {
				responses.InternalServerError(c, err)
				return
			}
		} else {
			if parentComment, err = findReplyParentComment(c, commentService, comment); err != nil {
				return
			}
			if err = commentService.DetrashCommentReply(comment, parentComment); err != nil {
				responses.InternalServerError(c, err)
				return
			}
		}

		responses.NoContent(c)
	}
}

func DeleteMyComment(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService  = service.New(c, ctx, dbConn)
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
		if comment, err = commentService.GetComment(bson.M{
			"$and": []bson.M{
				{"postauthoruid": bson.M{"$eq": me.UID}},
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
			if post, err = findCommentPost(c, commentService, comment); err != nil {
				return
			}
			if err = commentService.DeleteComment(comment, post); err != nil {
				responses.InternalServerError(c, err)
				return
			}
		} else {
			if parentComment, err = findReplyParentComment(c, commentService, comment); err != nil {
				return
			}
			if err = commentService.DeleteCommentReply(comment, parentComment); err != nil {
				responses.InternalServerError(c, err)
				return
			}
		}

		responses.NoContent(c)
	}
}

func findCommentPost(
	c *gin.Context,
	commentService *service.Service,
	comment *models.CommentModel,
) (post *models.PostModel, err error) {

	if post, err = commentService.GetPost(bson.M{
		"$and": []bson.M{
			{"deletedat": bson.M{"$eq": primitive.Null{}}},
			{"_id": bson.M{"$eq": comment.PostUid}}},
	}); err != nil {
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
	commentService *service.Service,
	reply *models.CommentModel,
) (comment *models.CommentModel, err error) {
	if comment, err = commentService.GetComment(bson.M{
		"$and": []bson.M{
			{"deletedat": bson.M{"$eq": primitive.Null{}}},
			{"_id": bson.M{"$eq": reply.ParentCommentUid}}},
	}); err != nil {
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
