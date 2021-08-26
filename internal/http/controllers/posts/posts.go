package posts

import (
	"context"
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
			dbConn *mongo.Database
			post   *models.PostModel
			postId primitive.ObjectID
			err    error
		)

		postQuery := c.Param("post")
		if postId, err = primitive.ObjectIDFromHex(postQuery); err != nil {
			postId = primitive.ObjectID{}
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if post, err = repositories.GetPost(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"deletedat": primitive.Null{}},
				bson.M{"publishedat": bson.M{"$ne": primitive.Null{}}},
				bson.M{"$or": []interface{}{
					bson.M{"_id": postId},
					bson.M{"slug": postQuery},
				}},
			}}); err != nil {
			responses.NotFound(c, err)
			return
		}

		responses.PublicPost(c, post)
	}
}

func GetPublicPosts(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			posts  []*models.PostModel
			err    error
		)

		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if posts, err = repositories.GetPosts(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"deletedat": primitive.Null{}},
				bson.M{"publishedat": bson.M{"$ne": primitive.Null{}}},
			}},
			helpers.GetShowQuery(c),
			helpers.GetOrderQuery(c),
			helpers.GetAscQuery(c)); err != nil {
			responses.NotFound(c, err)
			return
		}

		responses.PublicPosts(c, posts)
	}
}
