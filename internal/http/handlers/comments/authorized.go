package comments

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

func GetComment(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService  = service.New(c, ctx, dbConn)
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
		if comment, err = commentService.GetComment(bson.M{
			"_id": bson.M{"$eq": commentUid},
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

func GetComments(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel    = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService = service.New(c, ctx, dbConn)
			comments       []*models.CommentModel
			extraQuery     = []bson.M{}
			err            error
		)

		defer cancel()
		if trashQuery := c.DefaultQuery("trash", "false"); trashQuery == "true" {
			extraQuery = append(extraQuery,
				bson.M{"deletedat": bson.M{"$ne": primitive.Null{}}})
		} else {
			extraQuery = append(extraQuery,
				bson.M{"deletedat": bson.M{"$eq": primitive.Null{}}})
		}
		if comments, err = commentService.GetComments(bson.M{
			"$and": extraQuery,
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

func GetPostComments(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel    = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService = service.New(c, ctx, dbConn)
			comments       []*models.CommentModel
			postUid        primitive.ObjectID
			postUidParam   = c.Param("post")
			extraQuery     = []bson.M{}
			err            error
		)

		defer cancel()
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if trashQuery := c.DefaultQuery("trash", "false"); trashQuery == "true" {
			extraQuery = append(extraQuery,
				bson.M{"deletedat": bson.M{"$ne": primitive.Null{}}})
		} else {
			extraQuery = append(extraQuery,
				bson.M{"deletedat": bson.M{"$eq": primitive.Null{}}})
		}
		if comments, err = commentService.GetComments(bson.M{
			"$and": append(extraQuery,
				bson.M{"postuid": bson.M{"$eq": postUid}}),
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

func TrashComment(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel      = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService   = service.New(c, ctx, dbConn)
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
		if comment, err = commentService.GetComment(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
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

func DetrashComment(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService  = service.New(c, ctx, dbConn)
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
		if comment, err = commentService.GetComment(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
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

func DeleteComment(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService  = service.New(c, ctx, dbConn)
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
		if comment, err = commentService.GetComment(bson.M{
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
