package forms

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

type UpdateCategoryForm struct {
	Slug string `json:"slug" binding:"omitempty,max=100"`
	Name string `json:"name" binding:"omitempty,max=100"`
}

func (form *UpdateCategoryForm) Validate(
	svc *service.Service,
	ctx context.Context,
	target *models.CategoryModel,
) (err error) {
	if len(form.Slug) > 0 {
		if err = checkUpdateCategorySlug(svc, ctx, form.Slug, target); err != nil {
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

func checkUpdateCategorySlug(
	svc *service.Service,
	ctx context.Context,
	formSlug string,
	target *models.CategoryModel,
) (err error) {
	var categories []*models.CategoryModel

	if categories, err = svc.Category.GetMany(ctx, bson.M{
		"$and": []bson.M{
			{"_id": bson.M{"$ne": target.UID}},
			{"slug": bson.M{"$eq": formSlug}}},
	}); err != nil {
		return err
	}
	if len(categories) > 0 {
		return errors.New("slug exists")
	}

	return nil
}
