package posts

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/http/controllers/helpers"
	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

func GetMyPost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn   *mongo.Database
			me       *models.UserModel
			postData *models.PostModel
			postId   primitive.ObjectID
			err      error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}
		postIdQuery := c.Param("post")
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": "incorrent post id format"})
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if postData, err = repositories.GetPost(ctx, dbConn, bson.M{"$and": []interface{}{
			bson.M{"deletedAt": bson.M{"$exists": false}},
			bson.M{"author.username": me.Username},
			bson.M{"_id": postId},
		}}); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": "post not found"})
			return
		}

		responses.MyPost(c, postData)
	}
}

func GetMyPosts(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn     *mongo.Database
			me         *models.UserModel
			postsData  []*models.PostModel
			trashQuery interface{} = primitive.Null{}
			err        error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}
		if trashParam := c.DefaultQuery("trash", "false"); trashParam == "true" {
			trashQuery = bson.M{"$ne": primitive.Null{}}
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if postsData, err = repositories.GetPosts(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"deletedat": trashQuery},
				bson.M{"author.username": me.Username},
			}},
			helpers.GetShowQuery(c),
			helpers.GetOrderQuery(c),
			helpers.GetAscQuery(c)); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		responses.MyPosts(c, postsData)
	}
}

func CreatePost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			me     *models.UserModel
			post   *models.PostModel
			form   *forms.CreatePostForm
			err    error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}
		if form, err = requests.GetCreatePostForm(c); err != nil {
			responses.Basic(c, http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		post = forms.CreatePostModel(form, me)
		if err = repositories.CreatePost(ctx, dbConn, post); err != nil {
			var writeErr mongo.WriteException
			if errors.As(err, &writeErr) {
				responses.Basic(c, http.StatusUnprocessableEntity, gin.H{"message": writeErr.WriteErrors.Error()})
				return
			}
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responses.MyPost(c, post)
	}
}

func PublishMyPost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			me     *models.UserModel
			post   *models.PostModel
			postId primitive.ObjectID
			err    error
		)
		postIdQuery := c.Param("post")

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": "incorrent post id format"})
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if post, err = repositories.GetPost(ctx, dbConn, bson.M{"_id": postId}); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": "post not found for id: " + postIdQuery})
			return
		}
		if post.Author.Username != me.Username {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "unauthorized action"})
			return
		}
		if post.PublishedAt != nil {
			responses.Basic(c, http.StatusNoContent, nil)
			return
		}
		if err = repositories.PublishPost(ctx, dbConn, post); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responses.Basic(c, http.StatusNoContent, nil)
	}
}

func DepublishMyPost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			me     *models.UserModel
			post   *models.PostModel
			postId primitive.ObjectID
			err    error
		)
		postIdQuery := c.Param("post")

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": "incorrent post id format"})
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if post, err = repositories.GetPost(ctx, dbConn, bson.M{"_id": postId}); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": "post not found for id: " + postIdQuery})
			return
		}
		if post.Author.Username != me.Username {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "unauthorized action"})
			return
		}
		if post.PublishedAt == nil {
			responses.Basic(c, http.StatusNoContent, nil)
			return
		}
		if err = repositories.DepublishPost(ctx, dbConn, post); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responses.Basic(c, http.StatusNoContent, nil)
	}
}

func UpdateMyPost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			me     *models.UserModel
			post   *models.PostModel
			postId primitive.ObjectID
			form   *forms.UpdatePostForm
			err    error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}
		postIdQuery := c.Param("post")
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": "incorrent post id format"})
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if post, err = repositories.GetPost(ctx, dbConn, bson.M{"_id": postId}); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": "post not found for id: " + postIdQuery})
			return
		}
		if post.Author.Username != me.Username {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "unauthorized action"})
			return
		}
		if form, err = requests.GetUpdatePostForm(c); err != nil {
			responses.Basic(c, http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}
		if err = repositories.UpdatePost(ctx, dbConn, forms.UpdatePostModel(form, post)); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responses.Basic(c, http.StatusNoContent, nil)
	}
}

func TrashMyPost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			me     *models.UserModel
			post   *models.PostModel
			postId primitive.ObjectID
			err    error
		)
		postIdQuery := c.Param("post")

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": "incorrent post id format"})
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if post, err = repositories.GetPost(ctx, dbConn, bson.M{"_id": postId}); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": "post not found for id: " + postIdQuery})
			return
		}
		if post.Author.Username != me.Username {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "unauthorized action"})
			return
		}
		if post.DeletedAt != nil {
			responses.Basic(c, http.StatusNoContent, nil)
			return
		}
		if err = repositories.TrashPost(ctx, dbConn, post); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responses.Basic(c, http.StatusNoContent, nil)
	}
}

func DetrashMyPost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			me     *models.UserModel
			post   *models.PostModel
			postId primitive.ObjectID
			err    error
		)
		postIdQuery := c.Param("post")

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": "incorrent post id format"})
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if post, err = repositories.GetPost(ctx, dbConn, bson.M{"_id": postId}); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": "post not found for id: " + postIdQuery})
			return
		}
		if post.Author.Username != me.Username {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "unauthorized action"})
			return
		}
		if post.DeletedAt == nil {
			responses.Basic(c, http.StatusNoContent, nil)
			return
		}
		if err = repositories.DetrashPost(ctx, dbConn, post); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responses.Basic(c, http.StatusNoContent, nil)
	}
}

func DeleteMyPost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			me     *models.UserModel
			post   *models.PostModel
			postId primitive.ObjectID
			err    error
		)
		postIdQuery := c.Param("post")

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": "incorrent post id format"})
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if post, err = repositories.GetPost(ctx, dbConn, bson.M{"_id": postId}); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": "post not found for id: " + postIdQuery})
			return
		}
		if post.Author.Username != me.Username {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "unauthorized action"})
			return
		}
		if err = repositories.DeletePost(ctx, dbConn, post); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responses.Basic(c, http.StatusNoContent, nil)
	}
}

func GetMyPostComment(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{})
	}
}

func GetMyPostComments(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{})
	}
}

func TrashMyPostComment(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{})
	}
}

func DeleteMyPostComment(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{})
	}
}
