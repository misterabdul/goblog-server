package posts

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/http/controllers/helpers"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

func GetPublicPost(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn   *mongo.Database
			postData *models.PostModel
			postId   primitive.ObjectID
			err      error
		)
		postIdQuery := c.Param("post")

		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": "post not found"})
			return
		}
		if postData, err = repositories.GetPost(ctx, dbConn, bson.M{"$and": []interface{}{
			bson.M{"deletedAt": bson.M{"$exists": false}},
			bson.M{"publishedAt": bson.M{"$exists": true}},
			bson.M{"_id": postId},
		}}); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": "post not found"})
			return
		}

		responses.PublicPost(c, postData)
	}
}

func GetPublicPostSlug(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn   *mongo.Database
			postData *models.PostModel
			err      error
		)
		postSlugQuery := c.Param("post")

		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if postData, err = repositories.GetPost(ctx, dbConn, bson.M{"$and": []interface{}{
			bson.M{"deletedAt": bson.M{"$exists": false}},
			bson.M{"publishedAt": bson.M{"$exists": true}},
			bson.M{"slug": postSlugQuery},
		}}); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": "post not found"})
			return
		}

		responses.PublicPost(c, postData)
	}
}

func GetPublicPosts(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var dbConn *mongo.Database
		var postsData []*models.PostModel
		var err error

		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if postsData, err = repositories.GetPosts(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"deletedAt": bson.M{"$exists": false}},
				bson.M{"publishedAt": bson.M{"$exists": true}},
			}},
			helpers.GetShowQuery(c),
			helpers.GetOrderQuery(c),
			helpers.GetAscQuery(c)); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		responses.PublicPosts(c, postsData)
	}
}
