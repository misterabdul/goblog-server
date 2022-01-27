package posts

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

func GetPublicPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			postService = service.New(c, ctx, dbConn)
			post        *models.PostModel
			postContent *models.PostContentModel
			postUid     interface{}
			postParam   = c.Param("post")
			err         error
		)

		defer cancel()
		if postUid, err = primitive.ObjectIDFromHex(postParam); err != nil {
			postUid = nil
		}
		if post, postContent, err = postService.GetPostWithContent(bson.M{
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

		responses.PublicPost(c, post, postContent)
	}
}

func GetPublicPosts(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			postService = service.New(c, ctx, dbConn)
			posts       []*models.PostModel
			err         error
		)

		defer cancel()
		if posts, err = postService.GetPosts(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}}},
		}); err != nil {
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
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			postService = service.New(c, ctx, dbConn)
			searchQuery = c.Query("q")
			posts       []*models.PostModel
			err         error
		)

		defer cancel()
		if posts, err = postService.GetPosts(bson.M{
			"$text": bson.M{"$search": searchQuery},
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}}},
		}); err != nil {
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
