package forms

import (
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

type UpdateCategoryForm struct {
	Slug string `json:"slug" binding:"omitempty,max=100"`
	Name string `json:"name" binding:"omitempty,max=100"`
}

func (form *UpdateCategoryForm) Validate(
	categoryService *service.Service,
) (err error) {
	if len(form.Slug) > 0 {
		if err = checkCategorySlug(categoryService, form.Slug); err != nil {
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
