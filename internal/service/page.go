package service

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/database/repositories"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	customMongo "github.com/misterabdul/goblog-server/pkg/mongo"
)

type PageService struct {
	c                 *gin.Context
	ctx               context.Context
	dbConn            *mongo.Database
	repository        *repositories.PageRepository
	contentRepository *repositories.PageContentRepository
}

func NewPageService(
	c *gin.Context,
	ctx context.Context,
	dbConn *mongo.Database,
) *PageService {

	return &PageService{
		c:                 c,
		ctx:               ctx,
		dbConn:            dbConn,
		repository:        repositories.NewPageRepository(dbConn),
		contentRepository: repositories.NewPageContentRepository(dbConn)}
}

// Get single page
func (s *PageService) GetPage(filter interface{}) (
	page *models.PageModel,
	err error,
) {

	return s.repository.ReadOne(
		s.ctx, filter)
}

// Get single page with its content
func (s *PageService) GetPageWithContent(filter interface{}) (
	page *models.PageModel,
	content *models.PageContentModel,
	err error,
) {
	if page, err = s.repository.ReadOne(
		s.ctx, filter,
	); err != nil {
		return nil, nil, err
	}
	if page == nil {
		return nil, nil, nil
	}
	if content, err = s.contentRepository.ReadOne(
		s.ctx, bson.M{
			"_id": bson.M{"$eq": page.UID}},
	); err != nil {
		return page, nil, err
	}

	return page, content, nil
}

// Get multiple pages
func (s *PageService) GetPages(filter interface{}) (
	pages []*models.PageModel,
	err error,
) {

	return s.repository.ReadMany(
		s.ctx, filter,
		internalGin.GetFindOptions(s.c))
}

// Get total pages count
func (s *PageService) GetPageCount(filter interface{}) (
	count int64, err error,
) {

	return s.repository.Count(
		s.ctx, filter,
		internalGin.GetCountOptions(s.c))
}

// Create new page with its content
func (s *PageService) CreatePage(
	page *models.PageModel,
	content *models.PageContentModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	page.UID = primitive.NewObjectID()
	page.CreatedAt = now
	page.UpdatedAt = now
	page.DeletedAt = nil
	content.UID = page.UID

	return customMongo.Transaction(s.ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = s.repository.Save(
				sCtx, page,
			); sErr != nil {
				return sErr
			}
			if sErr = s.contentRepository.Save(
				sCtx, content,
			); sErr != nil {
				return sErr
			}

			return nil
		})
}

// Mark the page published
func (s *PageService) PublishPage(
	page *models.PageModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	page.PublishedAt = now

	return s.repository.Update(
		s.ctx, page)
}

// Remove published mark from the page
func (s *PageService) DepublishPage(
	page *models.PageModel,
) (err error) {
	page.PublishedAt = nil

	return s.repository.Update(
		s.ctx, page)
}

// Update page
func (s *PageService) UpdatePage(
	page *models.PageModel,
	content *models.PageContentModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	page.UpdatedAt = now

	return customMongo.Transaction(s.ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = s.repository.Update(
				sCtx, page,
			); sErr != nil {
				return sErr
			}
			if sErr = s.contentRepository.Update(
				sCtx, content,
			); sErr != nil {
				return sErr
			}

			return nil
		})
}

// Delete page to trash
func (s *PageService) TrashPage(
	page *models.PageModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	page.DeletedAt = now

	return s.repository.Update(
		s.ctx, page)
}

// Restore page from trash
func (s *PageService) DetrashPage(
	page *models.PageModel,
) (err error) {
	page.DeletedAt = nil

	return s.repository.Update(
		s.ctx, page)
}

// Permanently delete page
func (s *PageService) DeletePage(
	page *models.PageModel,
	content *models.PageContentModel,
) (err error) {

	return customMongo.Transaction(s.ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = s.repository.Delete(
				sCtx, page,
			); sErr != nil {
				return sErr
			}
			if sErr = s.contentRepository.Delete(
				sCtx, content,
			); sErr != nil {
				return sErr
			}

			return nil
		})
}
