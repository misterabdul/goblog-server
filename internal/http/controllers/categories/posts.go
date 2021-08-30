package categories

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/controllers/helpers"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

func GetPublicCategoryPosts(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			posts         []*models.PostModel
			categoryId    primitive.ObjectID
			categoryQuery = c.Param("category")
			err           error
		)

		if categoryId, err = primitive.ObjectIDFromHex(categoryQuery); err != nil {
			categoryId = primitive.ObjectID{}
		}
		if posts, err = repositories.GetPosts(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}},
				{"$or": []bson.M{
					{"_id": categoryId},
					{"slug": categoryQuery},
				}},
			}},
			helpers.GetShowQuery(c),
			helpers.GetOrderQuery(c),
			helpers.GetAscQuery(c)); err != nil {
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
