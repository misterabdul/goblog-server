package posts

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

func GetPublicPostComment(
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
				{"deletedat": primitive.Null{}},
				{"_id": commentId}},
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
				{"deletedat": primitive.Null{}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": comment.PostUid}},
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
			postId         primitive.ObjectID
			postQuery      = c.Param("post")
			err            error
		)

		defer cancel()
		if postId, err = primitive.ObjectIDFromHex(postQuery); err != nil {
			postId = primitive.ObjectID{}
		}
		if post, err = commentService.GetPost(bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}},
				{"$or": []bson.M{
					{"_id": postId},
					{"slug": postQuery}}}},
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
				{"deletedat": primitive.Null{}},
				{"$or": []bson.M{
					{"postslug": postQuery},
					{"postuid": postId}}}},
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
