package forms

import (
	"context"

	"github.com/misterabdul/goblog-server/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UpdateCategoryForm struct {
	Slug string `json:"slug" binding:"omitempty,max=100"`
	Name string `json:"name" binding:"omitempty,max=100"`
}

func (form *UpdateCategoryForm) Validate(
	ctx context.Context,
	dbConn *mongo.Database,
) (err error) {
	if len(form.Slug) > 0 {
		if err = checkCategorySlug(ctx, dbConn, form.Slug); err != nil {
			return err
		}
	}

	return nil
}

func (form *UpdateCategoryForm) ToCategoryModel(
	category *models.CategoryModel,
) (updatedCategory *models.CategoryModel) {
	if len(form.Slug) > 0 {
		category.Slug = form.Slug
	}
	if len(form.Name) > 0 {
		category.Name = form.Name
	}

	return category
}
