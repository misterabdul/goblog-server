package categories

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/service"
)

// @Tags        Category (Editor)
// @Summary     Get Category
// @Description Get category.
// @Router      /v1/auth/editor/category/{uid} [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Category's UID or slug"
// @Success     200 {object} object{data=object{uid=string,slug=string,name=string,updatedAt=time,createdAt=time}}
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func GetCategory(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel      = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService  = service.NewCategoryService(c, ctx, dbConn)
			category         *models.CategoryModel
			categoryUid      primitive.ObjectID
			categoryUidParam = c.Param("category")
			err              error
		)

		defer cancel()
		if categoryUid, err = primitive.ObjectIDFromHex(categoryUidParam); err != nil {
			responses.NotFound(c, err)
			return
		}
		if category, err = categoryService.GetCategory(bson.M{
			"_id": bson.M{"$eq": categoryUid},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if category == nil {
			responses.NotFound(c, errors.New("category not found"))
			return
		}

		responses.AuthorizedCategory(c, category)
	}
}

// @Tags        Category (Editor)
// @Summary     Get Categories
// @Description Get categories.
// @Router      /v1/auth/editor/categories [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       show  query    int     false "Number of data to be shown."
// @Param       page  query    int     false "Selected page of data."
// @Param       order query    string  false "Selected field to order data with."
// @Param       asc   query    boolean false "Ascending or descending."
// @Param       type  query    string  false "Filter data by type, e.g.: ?type=trash."
// @Success     200   {object} object{data=[]object{uid=string,slug=string,name=string,updatedAt=time,createdAt=time}}
// @Success     204
// @Failure     401   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func GetCategories(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService = service.NewCategoryService(c, ctx, dbConn)
			categories      []*models.CategoryModel
			queryParams     = readCommonQueryParams(c)
			err             error
		)

		defer cancel()
		if categories, err = categoryService.GetCategories(bson.M{
			"$and": queryParams,
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(categories) == 0 {
			responses.NoContent(c)
			return
		}

		responses.AuthorizedCategories(c, categories)
	}
}

// @Tags        Category (Editor)
// @Summary     Get Categories Stats
// @Description Get categories's stats.
// @Router      /v1/auth/editor/categories/stats [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       show  query    int     false "Number of data to be shown."
// @Param       page  query    int     false "Selected page of data."
// @Param       order query    string  false "Selected field to order data with."
// @Param       asc   query    boolean false "Ascending or descending."
// @Param       type  query    string  false "Filter data by type, e.g.: ?type=trash."
// @Success     200   {object} object{data=object{currentPage=int,totalPages=int,itemsPerPage=int,totalItems=int}}
// @Failure     401   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func GetCategoriesStats(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService = service.NewCategoryService(c, ctx, dbConn)
			count           int64
			queryParams     = readCommonQueryParams(c)
			err             error
		)

		defer cancel()
		if count, err = categoryService.GetCategoryCount(bson.M{
			"$and": queryParams,
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.ResourceStats(c, count)
	}
}

// @Tags        Category (Editor)
// @Summary     Create Category
// @Description Create a new category.
// @Router      /v1/auth/editor/category [post]
// @Security    BearerAuth
// @Accept      application/json
// @Accept      application/msgpack
// @Produce     application/json
// @Produce     application/msgpack
// @Param       form body     object{slug=string,name=string} true "Create category form"
// @Success     200  {object} object{data=object{uid=string,slug=string,name=string,updatedAt=time,createdAt=time}}
// @Failure     401  {object} object{message=string}
// @Failure     422  {object} object{message=string}
// @Failure     500  {object} object{message=string}
func CreateCategory(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService = service.NewCategoryService(c, ctx, dbConn)
			category        *models.CategoryModel
			form            *forms.CreateCategoryForm
			err             error
		)

		defer cancel()
		if form, err = requests.GetCreateCategoryForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if err = form.Validate(categoryService); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		category = form.ToCategoryModel()
		if err = categoryService.CreateCategory(category); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.AuthorizedCategory(c, category)
	}
}

// @Tags        Category (Editor)
// @Summary     Update Category
// @Description Update a category.
// @Router      /v1/auth/editor/category/{uid} [put]
// @Router      /v1/auth/editor/category/{uid} [patch]
// @Security    BearerAuth
// @Accept      application/json
// @Accept      application/msgpack
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid  path     string                          true "Category's UID or slug"
// @Param       form body     object{slug=string,name=string} true "Create category form"
// @Success     204
// @Failure     401  {object} object{message=string}
// @Failure     404  {object} object{message=string}
// @Failure     422  {object} object{message=string}
// @Failure     500  {object} object{message=string}
func UpdateCategory(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel      = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService  = service.NewCategoryService(c, ctx, dbConn)
			category         *models.CategoryModel
			updatedCategory  *models.CategoryModel
			categoryUid      primitive.ObjectID
			categoryUidParam = c.Param("category")
			form             *forms.UpdateCategoryForm
			err              error
		)

		defer cancel()
		if categoryUid, err = primitive.ObjectIDFromHex(categoryUidParam); err != nil {
			responses.IncorrectCategoryId(c, err)
			return
		}
		if category, err = categoryService.GetCategory(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": categoryUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if category == nil {
			responses.NotFound(c, errors.New("category not found"))
			return
		}
		if form, err = requests.GetUpdateCategoryForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if err = form.Validate(categoryService, category); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		updatedCategory = form.ToCategoryModel(category)
		if err = categoryService.UpdateCategory(updatedCategory); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Category (Editor)
// @Summary     Delete Category (Soft)
// @Description Delete a category (soft-deleted).
// @Router      /v1/auth/editor/category/{uid} [delete]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Category's UID or slug"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func TrashCategory(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel      = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService  = service.NewCategoryService(c, ctx, dbConn)
			category         *models.CategoryModel
			categoryUid      primitive.ObjectID
			categoryUidParam = c.Param("category")
			err              error
		)

		defer cancel()
		if categoryUid, err = primitive.ObjectIDFromHex(categoryUidParam); err != nil {
			responses.IncorrectCategoryId(c, err)
			return
		}
		if category, err = categoryService.GetCategory(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": categoryUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if category == nil {
			responses.NotFound(c, errors.New("category not found"))
			return
		}
		if category.DeletedAt != nil {
			responses.NoContent(c)
			return
		}
		if err = categoryService.TrashCategory(category); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Category (Editor)
// @Summary     Detrash Category
// @Description Restore a deleted category (soft-deleted).
// @Router      /v1/auth/editor/category/{uid}/detrash [put]
// @Router      /v1/auth/editor/category/{uid}/detrash [patch]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Category's UID or slug"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DetrashCategory(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel      = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService  = service.NewCategoryService(c, ctx, dbConn)
			category         *models.CategoryModel
			categoryUid      primitive.ObjectID
			categoryUidParam = c.Param("category")
			err              error
		)

		defer cancel()
		if categoryUid, err = primitive.ObjectIDFromHex(categoryUidParam); err != nil {
			responses.IncorrectCategoryId(c, err)
			return
		}
		if category, err = categoryService.GetCategory(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": bson.M{"$eq": categoryUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if category == nil {
			responses.NotFound(c, errors.New("category not found"))
			return
		}
		if category.DeletedAt == nil {
			responses.NoContent(c)
			return
		}
		if err = categoryService.DetrashCategory(category); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Category (Editor)
// @Summary     Delete Category (Permanent)
// @Description Delete a category (permanent).
// @Router      /v1/auth/editor/category/{uid}/permanent [delete]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Category's UID or slug"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DeleteCategory(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel      = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService  = service.NewCategoryService(c, ctx, dbConn)
			category         *models.CategoryModel
			categoryUid      primitive.ObjectID
			categoryUidParam = c.Param("category")
			err              error
		)

		defer cancel()
		if categoryUid, err = primitive.ObjectIDFromHex(categoryUidParam); err != nil {
			responses.IncorrectCategoryId(c, err)
			return
		}
		if category, err = categoryService.GetCategory(bson.M{
			"_id": bson.M{"$eq": categoryUid},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if category == nil {
			responses.NotFound(c, errors.New("category not found"))
			return
		}
		if err = categoryService.DeleteCategory(category); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func readCommonQueryParams(c *gin.Context) []bson.M {
	var (
		typeParam  = c.DefaultQuery("type", "active")
		extraQuery = []bson.M{}
	)

	switch true {
	case typeParam == "trash":
		extraQuery = append(extraQuery,
			bson.M{"deletedat": bson.M{"$ne": primitive.Null{}}})
	case typeParam == "active":
		fallthrough
	default:
		extraQuery = append(extraQuery,
			bson.M{"deletedat": bson.M{"$eq": primitive.Null{}}})
	}

	return extraQuery
}
