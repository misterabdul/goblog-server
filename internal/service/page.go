package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/database/repositories"
	customMongo "github.com/misterabdul/goblog-server/pkg/mongo"
)

type page struct {
	dbConn *mongo.Database
}

func newPageService(
	dbConn *mongo.Database,
) (service *page) {

	return &page{dbConn: dbConn}
}

// Get single page
func (s *page) GetOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (page *models.PageModel, err error) {

	return repositories.ReadOnePage(
		s.dbConn, ctx, filter, opts...)
}

// Get single page with its content
func (s *page) GetOneWithContent(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (page *models.PageModel, content *models.PageContentModel, err error) {
	if page, err = repositories.ReadOnePage(s.dbConn, ctx, filter, opts...); err != nil {
		return nil, nil, err
	}
	if page == nil {
		return nil, nil, nil
	}
	if content, err = repositories.ReadOnePageContent(s.dbConn, ctx,
		bson.M{"_id": bson.M{"$eq": page.UID}}, opts...,
	); err != nil {
		return page, nil, err
	}

	return page, content, nil
}

// Get multiple pages
func (s *page) GetMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (pages []*models.PageModel, err error) {

	return repositories.ReadManyPages(
		s.dbConn, ctx, filter, opts...)
}

// Get total pages count
func (s *page) Count(
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {

	return repositories.CountPages(
		s.dbConn, ctx, filter, opts...)
}

// Create new page with its content
func (s *page) SaveOneWithContent(
	ctx context.Context,
	page *models.PageModel,
	content *models.PageContentModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	page.UID = primitive.NewObjectID()
	page.CreatedAt = now
	page.UpdatedAt = now
	page.DeletedAt = nil
	content.UID = page.UID

	return customMongo.Transaction(ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = repositories.SaveOnePage(
				dbConn, sCtx, page, opts...,
			); sErr != nil {
				return sErr
			}
			if sErr = repositories.SaveOnePageContent(
				dbConn, sCtx, content, opts...,
			); sErr != nil {
				return sErr
			}

			return nil
		})
}

// Mark the page published
func (s *page) PublishOne(
	ctx context.Context,
	page *models.PageModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	page.PublishedAt = now

	return repositories.UpdateOnePage(
		s.dbConn, ctx, page, opts...)
}

// Remove published mark from the page
func (s *page) DepublishOne(
	ctx context.Context,
	page *models.PageModel,
	opts ...*options.UpdateOptions,
) (err error) {
	page.PublishedAt = nil

	return repositories.UpdateOnePage(
		s.dbConn, ctx, page, opts...)
}

// Update page
func (s *page) UpdateOneWithContent(
	ctx context.Context,
	page *models.PageModel,
	content *models.PageContentModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	page.UpdatedAt = now

	return customMongo.Transaction(ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = repositories.UpdateOnePage(
				dbConn, sCtx, page, opts...,
			); sErr != nil {
				return sErr
			}
			if sErr = repositories.UpdateOnePageContent(
				dbConn, sCtx, content, opts...,
			); sErr != nil {
				return sErr
			}

			return nil
		})
}

// Delete page to trash
func (s *page) TrashOne(
	ctx context.Context,
	page *models.PageModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	page.DeletedAt = now

	return repositories.UpdateOnePage(
		s.dbConn, ctx, page, opts...)
}

// Restore page from trash
func (s *page) RestoreOne(
	ctx context.Context,
	page *models.PageModel,
	opts ...*options.UpdateOptions,
) (err error) {
	page.DeletedAt = nil

	return repositories.UpdateOnePage(
		s.dbConn, ctx, page, opts...)
}

// Permanently delete page
func (s *page) DeleteOneWithContent(
	ctx context.Context,
	page *models.PageModel,
	content *models.PageContentModel,
	opts ...*options.DeleteOptions,
) (err error) {

	return customMongo.Transaction(ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = repositories.DeleteOnePage(
				dbConn, sCtx, page, opts...,
			); sErr != nil {
				return sErr
			}
			if sErr = repositories.DeleteOnePageContent(
				dbConn, sCtx, content, opts...,
			); sErr != nil {
				return sErr
			}

			return nil
		})
}
