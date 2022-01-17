package categories

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

func GetCategory(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService = service.New(c, ctx, dbConn)
			category        *models.CategoryModel
			categoryId      primitive.ObjectID
			categoryIdQuery = c.Param("category")
			err             error
		)

		defer cancel()
		if categoryId, err = primitive.ObjectIDFromHex(categoryIdQuery); err != nil {
			responses.NotFound(c, err)
			return
		}
		if category, err = categoryService.GetCategory(bson.M{
			"$and": []bson.M{
				{"_id": categoryId}},
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

func GetCategories(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService = service.New(c, ctx, dbConn)
			categories      []*models.CategoryModel
			trashQuery      interface{} = primitive.Null{}
			err             error
		)

		defer cancel()
		if trashParam := c.DefaultQuery("trash", "false"); trashParam == "true" {
			trashQuery = bson.M{"$ne": primitive.Null{}}
		}
		if categories, err = categoryService.GetCategories(bson.M{
			"$and": []bson.M{
				{"deletedat": trashQuery}},
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

func CreateCategory(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService = service.New(c, ctx, dbConn)
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

func UpdateCategory(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService = service.New(c, ctx, dbConn)
			category        *models.CategoryModel
			updatedCategory *models.CategoryModel
			categoryId      primitive.ObjectID
			form            *forms.UpdateCategoryForm
			categoryIdQuery = c.Param("category")
			err             error
		)

		defer cancel()
		if categoryId, err = primitive.ObjectIDFromHex(categoryIdQuery); err != nil {
			responses.IncorrectCategoryId(c, err)
			return
		}
		if category, err = categoryService.GetCategory(bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": categoryId}},
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
		if err = form.Validate(categoryService); err != nil {
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

func TrashCategory(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService = service.New(c, ctx, dbConn)
			category        *models.CategoryModel
			categoryId      primitive.ObjectID
			categoryIdQuery = c.Param("category")
			err             error
		)

		defer cancel()
		if categoryId, err = primitive.ObjectIDFromHex(categoryIdQuery); err != nil {
			responses.IncorrectCategoryId(c, err)
			return
		}
		if category, err = categoryService.GetCategory(bson.M{
			"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": categoryId}},
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

func DetrashCategory(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService = service.New(c, ctx, dbConn)
			category        *models.CategoryModel
			categoryId      primitive.ObjectID
			categoryIdQuery = c.Param("category")
			err             error
		)

		defer cancel()
		if categoryId, err = primitive.ObjectIDFromHex(categoryIdQuery); err != nil {
			responses.IncorrectCategoryId(c, err)
			return
		}
		if category, err = categoryService.GetCategory(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": categoryId}},
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

func DeleteCategory(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			categoryService = service.New(c, ctx, dbConn)
			category        *models.CategoryModel
			categoryId      primitive.ObjectID
			categoryIdQuery = c.Param("category")
			err             error
		)

		defer cancel()
		if categoryId, err = primitive.ObjectIDFromHex(categoryIdQuery); err != nil {
			responses.IncorrectCategoryId(c, err)
			return
		}
		if category, err = categoryService.GetCategory(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": categoryId}},
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
