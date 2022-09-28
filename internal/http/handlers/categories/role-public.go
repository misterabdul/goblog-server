package categories

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/service"
)

// @Tags        Category (Public)
// @Summary     Get Category
// @Description Get category.
// @Router      /v1/category/{uid} [get]
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Category's UID or slug"
// @Success     200 {object} object{data=object{uid=string,slug=string,name=string}}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func GetPublicCategory(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel      = context.WithTimeout(context.Background(), maxCtxDuration)
			category         *models.CategoryModel
			categoryUidParam = c.Param("category")
			categoryUid      interface{}
			err              error
		)

		defer cancel()
		if categoryUid, err = primitive.ObjectIDFromHex(categoryUidParam); err != nil {
			categoryUid = nil
		}
		if category, err = svc.Category.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"$or": []bson.M{
					{"_id": bson.M{"$eq": categoryUid}},
					{"slug": bson.M{"$eq": categoryUidParam}}}}}},
		); err != nil {
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

// @Tags        Category (Public)
// @Summary     Get Categories
// @Description Get categories.
// @Router      /v1/categories [get]
// @Produce     application/json
// @Produce     application/msgpack
// @Success     200   {object} object{data=[]object{uid=string,slug=string,name=string}}
// @Param       show  query    int     false "Number of data to be shown."
// @Param       page  query    int     false "Selected page of data."
// @Param       order query    string  false "Selected field to order data with."
// @Param       asc   query    boolean false "Ascending or descending."
// @Failure     204
// @Failure     500   {object} object{message=string}
func GetPublicCategories(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			categories  []*models.CategoryModel
			err         error
		)

		defer cancel()
		if categories, err = svc.Category.GetMany(ctx, bson.M{
			"deletedat": bson.M{"$eq": primitive.Null{}}},
			internalGin.GetFindOptions(c),
		); err != nil {
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
