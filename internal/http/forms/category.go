package forms

import (
	"github.com/misterabdul/goblog-server/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateCategoryForm struct {
	Slug string `json:"slug" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type UpdateCategoryForm struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

func CreateCategoryModel(form *CreateCategoryForm) (model *models.CategoryModel) {
	return &models.CategoryModel{
		UID:  primitive.NewObjectID(),
		Slug: form.Slug,
		Name: form.Name}
}

func UpdateCategoryModel(
	form *UpdateCategoryForm,
	category *models.CategoryModel,
) (model *models.CategoryModel) {
	if len(form.Slug) > 0 {
		category.Slug = form.Slug
	}
	if len(form.Name) > 0 {
		category.Name = form.Name
	}

	return category
}
