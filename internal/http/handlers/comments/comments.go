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
			ctx, cancel    = context.WithTimeout(context.Background(), maxCtxDuration)
			commentService = service.New(c, ctx, dbConn)
			comment        *models.CommentModel
			post           *models.PostModel
			commentId      primitive.ObjectID
			commentIdQuery = c.Param("comment")
			err            error
		)

		defer cancel()
		if commentId, err = primitive.ObjectIDFromHex(commentIdQuery); err != nil {
			responses.NotFound(c, err)
			return
		}
		if comment, err = commentService.GetComment(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": commentId}}},
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
			postId         interface{}
			postQuery      = c.Param("post")
			err            error
		)

		defer cancel()
		if postId, err = primitive.ObjectIDFromHex(postQuery); err != nil {
			postId = nil
		}
		if post, err = commentService.GetPost(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}},
				{"$or": []bson.M{
					{"_id": bson.M{"$eq": postId}},
					{"slug": bson.M{"$eq": postQuery}}}}},
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
		}); err != nil {
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
			commentId      interface{}
			commentQuery   = c.Param("comment")
			err            error
		)

		defer cancel()
		if commentId, err = primitive.ObjectIDFromHex(commentQuery); err != nil {
			responses.NotFound(c, errors.New("incorrent comment id format"))
			return
		}
		if comment, err = commentService.GetComment(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": commentId}}},
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
		}); err != nil {
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
			form           *forms.CreateCommentForm
			err            error
		)

		defer cancel()
		if form, err = requests.GetCreateCommentForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if err = form.Validate(commentService); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if comment, err = form.ToCommentModel(); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = commentService.CreateComment(comment); err != nil {
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
			form           *forms.CreateCommentReplyForm
			err            error
		)

		defer cancel()
		if form, err = requests.GetCreateCommentReplyForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if err = form.Validate(commentService); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if reply, err = form.ToCommentReplyModel(); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = commentService.CreateComment(reply); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.PublicComment(c, reply)
	}
}
