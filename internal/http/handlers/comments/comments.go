package comments

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

func GetPublicComment(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService  = service.New(c, ctx, dbConn)
			comment         *models.CommentModel
			post            *models.PostModel
			commentUid      primitive.ObjectID
			commentUidParam = c.Param("comment")
			err             error
		)

		defer cancel()
		if commentUid, err = primitive.ObjectIDFromHex(commentUidParam); err != nil {
			responses.NotFound(c, err)
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
		if post, err = commentService.GetPost(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": bson.M{"$eq": comment.PostUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
		}

		responses.PublicComment(c, comment)
	}
}

func GetPublicPostComments(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel    = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService = service.New(c, ctx, dbConn)
			comments       []*models.CommentModel
			post           *models.PostModel
			postUid        interface{}
			postParam      = c.Param("post")
			err            error
		)

		defer cancel()
		if postUid, err = primitive.ObjectIDFromHex(postParam); err != nil {
			postUid = nil
		}
		if post, err = commentService.GetPost(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}},
				{"$or": []bson.M{
					{"_id": bson.M{"$eq": postUid}},
					{"slug": bson.M{"$eq": postParam}}}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if comments, err = commentService.GetComments(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"parentcommentuid": bson.M{"$eq": primitive.Null{}}},
				{"postuid": bson.M{"$eq": post.UID}}},
		}, true); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(comments) == 0 {
			responses.NoContent(c)
			return
		}

		responses.PublicComments(c, comments)
	}
}

func GetPublicCommentReplies(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel    = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService = service.New(c, ctx, dbConn)
			replies        []*models.CommentModel
			comment        *models.CommentModel
			post           *models.PostModel
			commentUid     interface{}
			commentParam   = c.Param("comment")
			err            error
		)

		defer cancel()
		if commentUid, err = primitive.ObjectIDFromHex(commentParam); err != nil {
			responses.NotFound(c, errors.New("incorrent comment id format"))
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
		if post, err = commentService.GetPost(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": bson.M{"$eq": comment.PostUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if replies, err = commentService.GetComments(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"postuid": bson.M{"$eq": post.UID}},
				{"parentcommentuid": bson.M{"$eq": comment.UID}}},
		}, true); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(replies) == 0 {
			responses.NoContent(c)
			return
		}

		responses.PublicComments(c, replies)
	}
}

func CreatePublicPostComment(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel    = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService = service.New(c, ctx, dbConn)
			comment        *models.CommentModel
			post           *models.PostModel
			form           *forms.CreateCommentForm
			err            error
		)

		defer cancel()
		if form, err = requests.GetCreateCommentForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if post, err = form.Validate(commentService); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if comment, err = form.ToCommentModel(); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = commentService.CreateComment(comment, post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.PublicComment(c, comment)
	}
}

func CreatePublicCommentReply(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel    = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService = service.New(c, ctx, dbConn)
			reply          *models.CommentModel
			comment        *models.CommentModel
			form           *forms.CreateCommentReplyForm
			err            error
		)

		defer cancel()
		if form, err = requests.GetCreateCommentReplyForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if comment, err = form.Validate(commentService); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if reply, err = form.ToCommentReplyModel(); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = commentService.CreateCommentReply(reply, comment); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.PublicComment(c, reply)
	}
}
