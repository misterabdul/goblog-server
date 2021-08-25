package posts

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/http/controllers/helpers"
	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

func GetPost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			post   *models.PostModel
			postId primitive.ObjectID
			err    error
		)

		postIdQuery := c.Param("post")
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if post, err = repositories.GetPost(ctx, dbConn, bson.M{"$and": []interface{}{
			bson.M{"deletedat": primitive.Null{}},
			bson.M{"_id": postId},
		}}); err != nil {
			responses.NotFound(c, err)
			return
		}

		responses.AuthorizedPost(c, post)
	}
}

func GetPosts(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn     *mongo.Database
			posts      []*models.PostModel
			trashQuery interface{} = primitive.Null{}
			err        error
		)

		if trashParam := c.DefaultQuery("trash", "false"); trashParam == "true" {
			trashQuery = bson.M{"$ne": primitive.Null{}}
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if posts, err = repositories.GetPosts(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"deletedat": trashQuery},
			}},
			helpers.GetShowQuery(c),
			helpers.GetOrderQuery(c),
			helpers.GetAscQuery(c)); err != nil {
			responses.NotFound(c, err)
			return
		}

		responses.AuthorizedPosts(c, posts)
	}
}

func PublishPost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			post   *models.PostModel
			postId primitive.ObjectID
			err    error
		)

		postIdQuery := c.Param("post")
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if post, err = repositories.GetPost(ctx, dbConn, bson.M{"_id": postId}); err != nil {
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

func DepublishPost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			post   *models.PostModel
			postId primitive.ObjectID
			err    error
		)

		postIdQuery := c.Param("post")
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if post, err = repositories.GetPost(ctx, dbConn, bson.M{"_id": postId}); err != nil {
			responses.NotFound(c, err)
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

func UpdatePost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			post   *models.PostModel
			postId primitive.ObjectID
			form   *forms.UpdatePostForm
			err    error
		)

		postIdQuery := c.Param("post")
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if post, err = repositories.GetPost(ctx, dbConn, bson.M{"_id": postId}); err != nil {
			responses.NotFound(c, err)
			return
		}
		if form, err = requests.GetUpdatePostForm(c); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if err = repositories.UpdatePost(ctx, dbConn, forms.UpdatePostModel(form, post)); err != nil {
			var writeErr mongo.WriteException
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

func TrashPost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			post   *models.PostModel
			postId primitive.ObjectID
			err    error
		)

		postIdQuery := c.Param("post")
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if post, err = repositories.GetPost(ctx, dbConn, bson.M{"_id": postId}); err != nil {
			responses.NotFound(c, err)
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

func DetrashPost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			post   *models.PostModel
			postId primitive.ObjectID
			err    error
		)

		postIdQuery := c.Param("post")
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if post, err = repositories.GetPost(ctx, dbConn, bson.M{"_id": postId}); err != nil {
			responses.NotFound(c, err)
			return
		}
		if post.DeletedAt == nil {
			responses.NoContent(c)
			return
		}
		if err = repositories.DetrashPost(ctx, dbConn, post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DeletePost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			post   *models.PostModel
			postId primitive.ObjectID
			err    error
		)

		postIdQuery := c.Param("post")
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if post, err = repositories.GetPost(ctx, dbConn, bson.M{"_id": postId}); err != nil {
			responses.NotFound(c, err)
			return
		}
		if err = repositories.DeletePost(ctx, dbConn, post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func GetPostComment(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn    *mongo.Database
			comment   *models.CommentModel
			commentId primitive.ObjectID
			err       error
		)

		commentIdQuery := c.Param("comment")
		if commentId, err = primitive.ObjectIDFromHex(commentIdQuery); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if comment, err = repositories.GetComment(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"_id": commentId},
			}}); err != nil {
			responses.NotFound(c, err)
			return
		}

		responses.AuthorizedComment(c, comment)
	}
}

func GetPostComments(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn     *mongo.Database
			post       *models.PostModel
			comments   []*models.CommentModel
			postId     primitive.ObjectID
			trashQuery interface{} = primitive.Null{}
			err        error
		)

		postIdQuery := c.Param("post")
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if trashParam := c.DefaultQuery("trash", "false"); trashParam == "true" {
			trashQuery = bson.M{"$ne": primitive.Null{}}
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if post, err = repositories.GetPost(ctx, dbConn, bson.M{"_id": postId}); err != nil {
			responses.NotFound(c, err)
			return
		}
		if comments, err = repositories.GetComments(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"deletedat": trashQuery},
				bson.M{"_id": post.UID},
			}},
			helpers.GetShowQuery(c),
			helpers.GetOrderQuery(c),
			helpers.GetAscQuery(c)); err != nil {
			responses.NotFound(c, err)
			return
		}

		responses.AuthorizedComments(c, comments)
	}
}

func TrashPostComment(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn    *mongo.Database
			comment   *models.CommentModel
			commentId primitive.ObjectID
			err       error
		)

		commentIdQuery := c.Param("comment")
		if commentId, err = primitive.ObjectIDFromHex(commentIdQuery); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if comment, err = repositories.GetComment(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"_id": commentId},
			}}); err != nil {
			responses.NotFound(c, err)
			return
		}
		if _, err = repositories.GetPost(ctx, dbConn, bson.M{"_id": comment.PostUid}); err != nil {
			responses.NotFound(c, err)
			return
		}
		if comment.DeletedAt != nil {
			responses.NoContent(c)
			return
		}
		if err = repositories.TrashComment(ctx, dbConn, comment); err != nil {
			responses.InternalServerError(c, err)
		}

		responses.NoContent(c)
	}
}

func DetrashPostComment(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn    *mongo.Database
			comment   *models.CommentModel
			commentId primitive.ObjectID
			err       error
		)

		commentIdQuery := c.Param("comment")
		if commentId, err = primitive.ObjectIDFromHex(commentIdQuery); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if comment, err = repositories.GetComment(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"_id": commentId},
			}}); err != nil {
			responses.NotFound(c, err)
			return
		}
		if _, err = repositories.GetPost(ctx, dbConn, bson.M{"_id": comment.PostUid}); err != nil {
			responses.NotFound(c, err)
			return
		}
		if comment.DeletedAt == nil {
			responses.NoContent(c)
			return
		}
		if err = repositories.DetrashComment(ctx, dbConn, comment); err != nil {
			responses.InternalServerError(c, err)
		}

		responses.NoContent(c)
	}
}

func DeletePostComment(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn    *mongo.Database
			comment   *models.CommentModel
			commentId primitive.ObjectID
			err       error
		)

		commentIdQuery := c.Param("comment")
		if commentId, err = primitive.ObjectIDFromHex(commentIdQuery); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if comment, err = repositories.GetComment(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"_id": commentId},
			}}); err != nil {
			responses.NotFound(c, err)
			return
		}
		if _, err = repositories.GetPost(ctx, dbConn, bson.M{"_id": comment.PostUid}); err != nil {
			responses.NotFound(c, err)
			return
		}
		if err = repositories.DeleteComment(ctx, dbConn, comment); err != nil {
			responses.InternalServerError(c, err)
		}

		responses.NoContent(c)
	}
}
