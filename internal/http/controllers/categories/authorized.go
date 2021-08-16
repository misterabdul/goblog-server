package categories

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/http/controllers/helpers"
	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

func GetCategory(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn       *mongo.Database
			categoryData *models.CategoryModel
			categoryId   primitive.ObjectID
			err          error
		)
		categoryIdQuery := c.Param("category")

		if categoryId, err = primitive.ObjectIDFromHex(categoryIdQuery); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": "category not found"})
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if categoryData, err = repositories.GetCategory(ctx, dbConn, bson.M{"$and": []interface{}{
			bson.M{"_id": categoryId},
		}}); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": "category not found"})
			return
		}

		responses.AuthorizedCategory(c, categoryData)
	}
}

func GetCategories(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn         *mongo.Database
			categoriesData []*models.CategoryModel
			trashQuery     interface{} = primitive.Null{}
			err            error
		)

		if trashParam := c.DefaultQuery("trash", "false"); trashParam == "true" {
			trashQuery = bson.M{"$ne": primitive.Null{}}
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if categoriesData, err = repositories.GetCategories(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"deletedat": trashQuery},
			}},
			helpers.GetShowQuery(c),
			helpers.GetOrderQuery(c),
			helpers.GetAscQuery(c)); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		responses.AuthorizedCategories(c, categoriesData)
	}
}

func CreateCategory(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn   *mongo.Database
			category *models.CategoryModel
			form     *forms.CreateCategoryForm
			err      error
		)

		if form, err = requests.GetCreateCategoryForm(c); err != nil {
			responses.Basic(c, http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		category = forms.CreateCategoryModel(form)
		if err = repositories.CreateCategory(ctx, dbConn, category); err != nil {
			var writeErr mongo.WriteException
			if errors.As(err, &writeErr) {
				responses.Basic(c, http.StatusUnprocessableEntity, gin.H{"message": writeErr.WriteErrors.Error()})
				return
			}
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responses.AuthorizedCategory(c, category)
	}
}

func UpdateCategory(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn     *mongo.Database
			category   *models.CategoryModel
			categoryId primitive.ObjectID
			form       *forms.UpdateCategoryForm
			err        error
		)

		categoryIdQuery := c.Param("category")
		if categoryId, err = primitive.ObjectIDFromHex(categoryIdQuery); err != nil {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": "incorrent category id format"})
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if category, err = repositories.GetCategory(ctx, dbConn, bson.M{"$and": []interface{}{
			bson.M{"deletedat": primitive.Null{}},
			bson.M{"_id": categoryId},
		}}); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": "category not found for id: " + categoryIdQuery})
			return
		}
		if form, err = requests.GetUpdateCategoryForm(c); err != nil {
			responses.Basic(c, http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}
		if err = repositories.UpdateCategory(ctx, dbConn, forms.UpdateCategoryModel(form, category)); err != nil {
			var writeErr mongo.WriteException
			if errors.As(err, &writeErr) {
				responses.Basic(c, http.StatusUnprocessableEntity, gin.H{"message": writeErr.WriteErrors.Error()})
				return
			}
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responses.Basic(c, http.StatusNoContent, nil)
	}
}

func TrashCategory(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn     *mongo.Database
			category   *models.CategoryModel
			categoryId primitive.ObjectID
			err        error
		)

		categoryIdQuery := c.Param("category")
		if categoryId, err = primitive.ObjectIDFromHex(categoryIdQuery); err != nil {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": "incorrent category id format"})
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if category, err = repositories.GetCategory(ctx, dbConn, bson.M{"_id": categoryId}); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": "categeory not found for id: " + categoryIdQuery})
			return
		}
		if category.DeletedAt != nil {
			responses.Basic(c, http.StatusNoContent, nil)
			return
		}
		if err = repositories.TrashCategory(ctx, dbConn, category); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responses.Basic(c, http.StatusNoContent, nil)
	}
}

func DetrashCategory(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn     *mongo.Database
			category   *models.CategoryModel
			categoryId primitive.ObjectID
			err        error
		)

		categoryIdQuery := c.Param("category")
		if categoryId, err = primitive.ObjectIDFromHex(categoryIdQuery); err != nil {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": "incorrent category id format"})
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if category, err = repositories.GetCategory(ctx, dbConn, bson.M{"_id": categoryId}); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": "categeory not found for id: " + categoryIdQuery})
			return
		}
		if category.DeletedAt == nil {
			responses.Basic(c, http.StatusNoContent, nil)
			return
		}
		if err = repositories.DetrashCategory(ctx, dbConn, category); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responses.Basic(c, http.StatusNoContent, nil)
	}
}

func DeleteCategory(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn     *mongo.Database
			category   *models.CategoryModel
			categoryId primitive.ObjectID
			err        error
		)

		categoryIdQuery := c.Param("category")
		if categoryId, err = primitive.ObjectIDFromHex(categoryIdQuery); err != nil {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": "incorrent category id format"})
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if category, err = repositories.GetCategory(ctx, dbConn, bson.M{"_id": categoryId}); err != nil {
			responses.Basic(c, http.StatusNotFound, gin.H{"message": "post not found for id: " + categoryIdQuery})
			return
		}
		if err = repositories.DeleteCategory(ctx, dbConn, category); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responses.Basic(c, http.StatusNoContent, nil)
	}
}
