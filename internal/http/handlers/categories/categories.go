package categories

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

func GetPublicCategory(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService = service.New(c, ctx, dbConn)
			category        *models.CategoryModel
			categoryUid     interface{}
			categoryParam   = c.Param("category")
			err             error
		)

		defer cancel()
		if categoryUid, err = primitive.ObjectIDFromHex(categoryParam); err != nil {
			categoryUid = nil
		}
		if category, err = categoryService.GetCategory(bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"$or": []bson.M{
					{"_id": bson.M{"$eq": categoryUid}},
					{"slug": bson.M{"$eq": categoryParam}}}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if category == nil {
			responses.NotFound(c, errors.New("category not found"))
			return
		}

		responses.PublicCategory(c, category)
	}
}

func GetPublicCategories(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService = service.New(c, ctx, dbConn)
			categories      []*models.CategoryModel
			err             error
		)

		defer cancel()
		if categories, err = categoryService.GetCategories(bson.M{
			"deletedat": primitive.Null{},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(categories) == 0 {
			responses.NoContent(c)
			return
		}

		responses.PublicCategories(c, categories)
	}
}
