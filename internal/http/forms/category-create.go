package forms

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

type CreateCategoryForm struct {
	Slug string `json:"slug" binding:"required,alphanum,max=100"`
	Name string `json:"name" binding:"required,max=100"`
}

func (form *CreateCategoryForm) Validate(
	categoryService *service.CategoryService,
) (err error) {
	if err = checkCategorySlug(categoryService, form.Slug); err != nil {
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

func checkCategorySlug(
	categoryService *service.CategoryService, formSlug string,
) (err error) {
	var categories []*models.CategoryModel

	if categories, err = categoryService.GetCategories(bson.M{
		"slug": bson.M{"$eq": formSlug},
	}); err != nil {
		return err
	}
	if len(categories) > 0 {
		return errors.New("slug exists")
	}

	return nil
}
