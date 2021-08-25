package categories

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

func GetPublicCategory(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn       *mongo.Database
			categoryData *models.CategoryModel
			categoryId   primitive.ObjectID
			err          error
		)

		categoryQuery := c.Param("category")
		if categoryId, err = primitive.ObjectIDFromHex(categoryQuery); err != nil {
			categoryId = primitive.ObjectID{}
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if categoryData, err = repositories.GetCategory(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"deletedat": primitive.Null{}},
				bson.M{"$or": []interface{}{
					bson.M{"_id": categoryId},
					bson.M{"slug": categoryQuery},
				}},
			}}); err != nil {
			responses.NotFound(c, err)
			return
		}

		responses.PublicCategory(c, categoryData)
	}
}

func GetPublicCategories(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn     *mongo.Database
			categories []*models.CategoryModel
			err        error
		)

		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if categories, err = repositories.GetCategories(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"deletedat": primitive.Null{}},
			}},
			helpers.GetShowQuery(c),
			helpers.GetOrderQuery(c),
			helpers.GetAscQuery(c)); err != nil {
			responses.NotFound(c, err)
			return
		}

		responses.PublicCategories(c, categories)
	}
}
