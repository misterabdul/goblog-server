package service

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/database/repositories"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
)

type CategoryService struct {
	c          *gin.Context
	ctx        context.Context
	dbConn     *mongo.Database
	repository *repositories.CategoryRepository
}

func NewCategoryService(
	c *gin.Context,
	ctx context.Context,
	dbConn *mongo.Database,
) *CategoryService {

	return &CategoryService{
		c:          c,
		ctx:        ctx,
		dbConn:     dbConn,
		repository: repositories.NewCategoryRepository(dbConn)}
}

// Get single category
func (s *CategoryService) GetCategory(filter interface{}) (
	category *models.CategoryModel,
	err error,
) {

	return s.repository.ReadOne(
		s.ctx, filter)
}

// Get multiple categories
func (s *CategoryService) GetCategories(filter interface{}) (
	categories []*models.CategoryModel,
	err error,
) {

	return s.repository.ReadMany(
		s.ctx, filter,
		internalGin.GetFindOptions(s.c))
}

// Get total categories count
func (s *CategoryService) GetCategoryCount(filter interface{}) (
	count int64, err error,
) {

	return s.repository.Count(
		s.ctx, filter,
		internalGin.GetCountOptions(s.c))
}

// Create new category
func (s *CategoryService) CreateCategory(
	category *models.CategoryModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	category.UID = primitive.NewObjectID()
	category.CreatedAt = now
	category.UpdatedAt = now
	category.DeletedAt = nil

	return s.repository.Save(
		s.ctx, category)
}

// Update category
func (s *CategoryService) UpdateCategory(
	category *models.CategoryModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	category.UpdatedAt = now

	return s.repository.Update(
		s.ctx, category)
}

// Delete category to trash
func (s *CategoryService) TrashCategory(
	category *models.CategoryModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	category.DeletedAt = now

	return s.repository.Update(
		s.ctx, category)
}

// Restore category from trash
func (s *CategoryService) DetrashCategory(
	category *models.CategoryModel,
) (err error) {
	category.DeletedAt = nil

	return s.repository.Update(
		s.ctx, category)
}

// Permanently delete category
func (s *CategoryService) DeleteCategory(
	category *models.CategoryModel,
) (err error) {

	return s.repository.Delete(
		s.ctx, category)
}
