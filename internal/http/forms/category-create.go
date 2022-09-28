package forms

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

type CreateCategoryForm struct {
	Slug string `json:"slug" binding:"required,alphanum,max=100"`
	Name string `json:"name" binding:"required,max=100"`
}

func (form *CreateCategoryForm) Validate(
	svc *service.Service,
	ctx context.Context,
) (err error) {
	if err = checkCategorySlug(svc, ctx, form.Slug); err != nil {
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
	svc *service.Service,
	ctx context.Context,
	formSlug string,
) (err error) {
	var categories []*models.CategoryModel

	if categories, err = svc.Category.GetMany(ctx, bson.M{
		"slug": bson.M{"$eq": formSlug},
	}); err != nil {
		return err
	}
	if len(categories) > 0 {
		return errors.New("slug exists")
	}

	return nil
}
