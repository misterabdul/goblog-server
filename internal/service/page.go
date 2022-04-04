package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/models"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/repositories"
	customMongo "github.com/misterabdul/goblog-server/pkg/mongo"
)

// Get single page
func (service *Service) GetPage(filter interface{}) (
	page *models.PageModel,
	err error,
) {
	return repositories.GetPage(
		service.ctx,
		service.dbConn,
		filter)
}

// Get single page with its content
func (service *Service) GetPageWithContent(filter interface{}) (
	page *models.PageModel,
	content *models.PageContentModel,
	err error,
) {
	if page, err = repositories.GetPage(
		service.ctx, service.dbConn, filter,
	); err != nil {
		return nil, nil, err
	}
	if page == nil {
		return nil, nil, nil
	}
	if content, err = repositories.GetPageContent(
		service.ctx, service.dbConn, bson.M{
			"_id": bson.M{"$eq": page.UID}},
	); err != nil {
		return page, nil, err
	}

	return page, content, nil
}

// Get multiple pages
func (service *Service) GetPages(filter interface{}) (
	pages []*models.PageModel,
	err error,
) {

	return repositories.GetPages(
		service.ctx,
		service.dbConn,
		filter,
		internalGin.GetFindOptions(service.c))
}

// Create new page with its content
func (service *Service) CreatePage(
	page *models.PageModel,
	content *models.PageContentModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	page.UID = primitive.NewObjectID()
	page.CreatedAt = now
	page.UpdatedAt = now
	page.DeletedAt = nil
	content.UID = page.UID

	return customMongo.Transaction(service.ctx, service.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = repositories.SavePage(
				sCtx, dbConn, page,
			); sErr != nil {
				return sErr
			}
			if sErr = repositories.SavePageContent(
				sCtx, dbConn, content,
			); sErr != nil {
				return sErr
			}

			return nil
		})
}

// Mark the page published
func (service *Service) PublishPage(
	page *models.PageModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	page.PublishedAt = now

	return repositories.UpdatePage(
		service.ctx,
		service.dbConn,
		page)
}

// Remove published mark from the page
func (service *Service) DepublishPage(
	page *models.PageModel,
) (err error) {
	page.PublishedAt = nil

	return repositories.UpdatePage(
		service.ctx,
		service.dbConn,
		page)
}

// Update page
func (service *Service) UpdatePage(
	page *models.PageModel,
	content *models.PageContentModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	page.UpdatedAt = now

	return customMongo.Transaction(service.ctx, service.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = repositories.UpdatePage(
				sCtx, dbConn, page,
			); sErr != nil {
				return sErr
			}
			if sErr = repositories.UpdatePageContent(
				sCtx, dbConn, content,
			); sErr != nil {
				return sErr
			}

			return nil
		})
}

// Delete page to trash
func (service *Service) TrashPage(
	page *models.PageModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	page.DeletedAt = now

	return repositories.UpdatePage(
		service.ctx,
		service.dbConn,
		page)
}

// Restore page from trash
func (service *Service) DetrashPage(
	page *models.PageModel,
) (err error) {
	page.DeletedAt = nil

	return repositories.UpdatePage(
		service.ctx,
		service.dbConn,
		page)
}

// Permanently delete page
func (service *Service) DeletePage(
	page *models.PageModel,
	content *models.PageContentModel,
) (err error) {

	return customMongo.Transaction(service.ctx, service.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = repositories.DeletePage(
				sCtx, dbConn, page,
			); sErr != nil {
				return sErr
			}
			if sErr = repositories.DeletePageContent(
				sCtx, dbConn, content,
			); sErr != nil {
				return sErr
			}

			return nil
		})
}
