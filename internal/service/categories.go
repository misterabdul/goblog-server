package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/database/repositories"
)

type category struct {
	dbConn *mongo.Database
}

func newCategoryService(
	dbConn *mongo.Database,
) (service *category) {

	return &category{dbConn: dbConn}
}

// Get single category
func (s *category) GetOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (category *models.CategoryModel, err error) {

	return repositories.ReadOneCategory(
		s.dbConn, ctx, filter, opts...)
}

// Get multiple categories
func (s *category) GetMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (categories []*models.CategoryModel, err error) {

	return repositories.ReadManyCategories(
		s.dbConn, ctx, filter, opts...)
}

// Get total categories count
func (s *category) Count(
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {

	return repositories.CountCategories(
		s.dbConn, ctx, filter, opts...)
}

// Create new category
func (s *category) SaveOne(
	ctx context.Context,
	category *models.CategoryModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	category.UID = primitive.NewObjectID()
	category.CreatedAt = now
	category.UpdatedAt = now
	category.DeletedAt = nil

	return repositories.SaveOneCategory(
		s.dbConn, ctx, category, opts...)
}

// Update category
func (s *category) UpdateOne(
	ctx context.Context,
	category *models.CategoryModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	category.UpdatedAt = now

	return repositories.UpdateOneCategory(
		s.dbConn, ctx, category, opts...)
}

// Delete category to trash
func (s *category) TrashOne(
	ctx context.Context,
	category *models.CategoryModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	category.DeletedAt = now

	return repositories.UpdateOneCategory(
		s.dbConn, ctx, category, opts...)
}

// Restore category from trash
func (s *category) RestoreOne(
	ctx context.Context,
	category *models.CategoryModel,
	opts ...*options.UpdateOptions,
) (err error) {
	category.DeletedAt = nil

	return repositories.UpdateOneCategory(
		s.dbConn, ctx, category, opts...)
}

// Permanently delete category
func (s *category) DeleteOne(
	ctx context.Context,
	category *models.CategoryModel,
	opts ...*options.DeleteOptions,
) (err error) {

	return repositories.DeleteOneCategory(
		s.dbConn, ctx, category, opts...)
}
