package posts

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/handlers/helpers"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

func GetPublicPost(
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
			postQuery   = c.Param("post")
			err         error
		)

		if postId, err = primitive.ObjectIDFromHex(postQuery); err != nil {
			postId = primitive.ObjectID{}
		}
		if post, postContent, err = repositories.GetPostWithContent(ctx, dbConn, bson.M{
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

		responses.PublicPost(c, post, postContent)
	}
}

func GetPublicPosts(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			posts []*models.PostModel
			err   error
		)

		if posts, err = repositories.GetPosts(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}}},
		}, helpers.GetFindOptionsPost(c)); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(posts) == 0 {
			responses.NoContent(c)
			return
		}

		responses.PublicPosts(c, posts)
	}
}

func SearchPublicPosts(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			searchQuery = c.Query("q")
			posts       []*models.PostModel
			err         error
		)

		if posts, err = repositories.GetPosts(ctx, dbConn, bson.M{
			"$text": bson.M{
				"$search": searchQuery},
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}}},
		}, helpers.GetFindOptionsPost(c)); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(posts) == 0 {
			responses.NoContent(c)
			return
		}

		responses.PublicPosts(c, posts)
	}
}
