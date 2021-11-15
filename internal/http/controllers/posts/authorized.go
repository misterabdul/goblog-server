package posts

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/controllers/helpers"
	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

func GetPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			post        *models.PostModel
			postContent *models.PostContentModel
			postId      primitive.ObjectID
			postIdQuery = c.Param("post")
			err         error
		)

		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, postContent, err = repositories.GetPostWithContent(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": postId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}

		responses.AuthorizedPost(c, post, postContent)
	}
}

func GetPosts(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			posts      []*models.PostModel
			trashQuery interface{} = primitive.Null{}
			err        error
		)

		if trashParam := c.DefaultQuery("trash", "false"); trashParam == "true" {
			trashQuery = bson.M{"$ne": primitive.Null{}}
		}
		if posts, err = repositories.GetPosts(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": trashQuery}},
		}, helpers.GetFindOptionsPost(c)); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(posts) == 0 {
			responses.NoContent(c)
			return
		}

		responses.AuthorizedPosts(c, posts)
	}
}

func PublishPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			post        *models.PostModel
			postId      primitive.ObjectID
			postIdQuery = c.Param("post")
			err         error
		)

		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, err = repositories.GetPost(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": postId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, err)
			return
		}
		if post.PublishedAt != nil {
			responses.NoContent(c)
			return
		}
		if err = repositories.PublishPost(ctx, dbConn, post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DepublishPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			post        *models.PostModel
			postId      primitive.ObjectID
			postIdQuery = c.Param("post")
			err         error
		)

		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, err = repositories.GetPost(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": postId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if post.PublishedAt == nil {
			responses.NoContent(c)
			return
		}
		if err = repositories.DepublishPost(ctx, dbConn, post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func UpdatePost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			post        *models.PostModel
			postContent *models.PostContentModel
			postId      primitive.ObjectID
			postIdQuery = c.Param("post")
			form        *forms.UpdatePostForm
			err         error
			writeErr    mongo.WriteException
		)

		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, postContent, err = repositories.GetPostWithContent(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": postId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if form, err = requests.GetUpdatePostForm(c); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		post, postContent = forms.UpdatePostModel(form, post, postContent)
		if err = repositories.UpdatePost(ctx, dbConn, post, postContent); err != nil {
			if errors.As(err, &writeErr) {
				responses.FormIncorrect(c, err)
				return
			}
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func TrashPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			post        *models.PostModel
			postId      primitive.ObjectID
			postIdQuery = c.Param("post")
			err         error
		)

		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, err = repositories.GetPost(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": postId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if post.DeletedAt != nil {
			responses.NoContent(c)
			return
		}
		if err = repositories.TrashPost(ctx, dbConn, post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DetrashPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			post        *models.PostModel
			postId      primitive.ObjectID
			postIdQuery = c.Param("post")
			err         error
		)

		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, err = repositories.GetPost(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": postId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if err = repositories.DetrashPost(ctx, dbConn, post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DeletePost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			post        *models.PostModel
			postContent *models.PostContentModel
			postId      primitive.ObjectID
			postIdQuery = c.Param("post")
			err         error
		)

		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, postContent, err = repositories.GetPostWithContent(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": postId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if err = repositories.DeletePost(ctx, dbConn, post, postContent); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func GetPostComment(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			comment        *models.CommentModel
			post           *models.PostModel
			commentId      primitive.ObjectID
			commentIdQuery = c.Param("comment")
			err            error
		)

		if commentId, err = primitive.ObjectIDFromHex(commentIdQuery); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if comment, err = repositories.GetComment(ctx, dbConn, bson.M{
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
		if post, err = repositories.GetPost(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}

		responses.AuthorizedComment(c, comment)
	}
}

func GetPostComments(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			post        *models.PostModel
			comments    []*models.CommentModel
			postId      primitive.ObjectID
			postIdQuery             = c.Param("post")
			trashQuery  interface{} = primitive.Null{}
			err         error
		)

		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if trashParam := c.DefaultQuery("trash", "false"); trashParam == "true" {
			trashQuery = bson.M{"$ne": primitive.Null{}}
		}
		if post, err = repositories.GetPost(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": postId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if comments, err = repositories.GetComments(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": trashQuery},
				{"_id": post.UID}},
		}, helpers.GetFindOptions(c)); err != nil {
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

func TrashPostComment(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			comment        *models.CommentModel
			post           *models.PostModel
			commentId      primitive.ObjectID
			commentIdQuery = c.Param("comment")
			err            error
		)

		if commentId, err = primitive.ObjectIDFromHex(commentIdQuery); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if comment, err = repositories.GetComment(ctx, dbConn, bson.M{
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
		if post, err = repositories.GetPost(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": comment.PostUid}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if err = repositories.TrashComment(ctx, dbConn, comment); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DetrashPostComment(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			comment        *models.CommentModel
			post           *models.PostModel
			commentId      primitive.ObjectID
			commentIdQuery = c.Param("comment")
			err            error
		)

		if commentId, err = primitive.ObjectIDFromHex(commentIdQuery); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if comment, err = repositories.GetComment(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": commentId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if comment == nil {
			responses.NotFound(c, errors.New("comment not found"))
			return
		}
		if post, err = repositories.GetPost(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": comment.PostUid}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if err = repositories.DetrashComment(ctx, dbConn, comment); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DeletePostComment(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			comment        *models.CommentModel
			post           *models.PostModel
			commentId      primitive.ObjectID
			commentIdQuery = c.Param("comment")
			err            error
		)

		if commentId, err = primitive.ObjectIDFromHex(commentIdQuery); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if comment, err = repositories.GetComment(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": commentId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if comment == nil {
			responses.NotFound(c, errors.New("comment not found"))
			return
		}
		if post, err = repositories.GetPost(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": comment.PostUid}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if err = repositories.DeleteComment(ctx, dbConn, comment); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}
