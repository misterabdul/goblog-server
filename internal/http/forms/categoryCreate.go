package forms

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

type CreateCategoryForm struct {
	Slug string `json:"slug" binding:"required,alphanum,max=100"`
	Name string `json:"name" binding:"required,max=100"`
}

func (form *CreateCategoryForm) Validate(
	ctx context.Context,
	dbConn *mongo.Database,
) (err error) {
	if err = checkCategorySlug(ctx, dbConn, form.Slug); err != nil {
		return err
	}

	return nil
}

func (form *CreateCategoryForm) ToCategoryModel() (model *models.CategoryModel) {
	return &models.CategoryModel{
		UID:  primitive.NewObjectID(),
		Slug: form.Slug,
		Name: form.Name}
}

func checkCategorySlug(ctx context.Context, dbConn *mongo.Database, formSlug string) (err error) {
	var (
		categories []*models.CategoryModel
	)

	if categories, err = repositories.GetCategories(ctx, dbConn, bson.M{
		"slug": bson.M{"$eq": formSlug},
	}); err != nil {
		return err
	}
	if len(categories) > 0 {
		return errors.New("slug exists")
	}

	return nil
}
