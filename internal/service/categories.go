package service

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/models"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

// Get single category
func (service *Service) GetCategory(filter interface{}) (
	category *models.CategoryModel,
	err error,
) {

	return repositories.GetCategory(
		service.ctx,
		service.dbConn,
		filter)
}

// Get multiple categories
func (service *Service) GetCategories(filter interface{}) (
	categories []*models.CategoryModel,
	err error,
) {

	return repositories.GetCategories(
		service.ctx,
		service.dbConn,
		filter,
		internalGin.GetFindOptions(service.c))
}

// Create new category
func (service *Service) CreateCategory(
	category *models.CategoryModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	category.UID = primitive.NewObjectID()
	category.CreatedAt = now
	category.UpdatedAt = now
	category.DeletedAt = nil

	return repositories.SaveCategory(
		service.ctx,
		service.dbConn,
		category)
}

// Update category
func (service *Service) UpdateCategory(
	category *models.CategoryModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	category.UpdatedAt = now

	return repositories.UpdateCategory(
		service.ctx,
		service.dbConn,
		category)
}

// Delete category to trash
func (service *Service) TrashCategory(
	category *models.CategoryModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	category.DeletedAt = now

	return repositories.UpdateCategory(
		service.ctx,
		service.dbConn,
		category)
}

// Restore category from trash
func (service *Service) DetrashCategory(
	category *models.CategoryModel,
) (err error) {
	category.DeletedAt = nil

	return repositories.UpdateCategory(
		service.ctx,
		service.dbConn,
		category)
}

// Permanently delete category
func (service *Service) DeleteCategory(
	category *models.CategoryModel,
) (err error) {

	return repositories.DeleteCategory(
		service.ctx,
		service.dbConn,
		category)
}
