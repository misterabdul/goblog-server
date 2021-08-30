package categories

import (
	"context"
	"errors"
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

func GetPublicCategory(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			category      *models.CategoryModel
			categoryId    primitive.ObjectID
			categoryQuery = c.Param("category")
			err           error
		)

		if categoryId, err = primitive.ObjectIDFromHex(categoryQuery); err != nil {
			categoryId = primitive.ObjectID{}
		}
		if category, err = repositories.GetCategory(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"$or": []bson.M{
					{"_id": categoryId},
					{"slug": categoryQuery},
				}},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if category == nil {
			responses.NotFound(c, errors.New("category not found"))
		}

		responses.PublicCategory(c, category)
	}
}

func GetPublicCategories(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			categories []*models.CategoryModel
			err        error
		)

		if categories, err = repositories.GetCategories(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": primitive.Null{}},
			}},
			helpers.GetShowQuery(c),
			helpers.GetOrderQuery(c),
			helpers.GetAscQuery(c)); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(categories) == 0 {
			responses.NoContent(c)
		}

		responses.PublicCategories(c, categories)
	}
}
