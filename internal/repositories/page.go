package repositories

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/models"
)

type PageRepository struct {
	collection *mongo.Collection
}

type PageContentRepository struct {
	collection *mongo.Collection
}

func NewPageRepository(
	dbConn *mongo.Database,
) *PageRepository {

	return &PageRepository{
		collection: dbConn.Collection("pages")}
}

func NewPageContentRepository(
	dbConn *mongo.Database,
) *PageContentRepository {

	return &PageContentRepository{
		collection: dbConn.Collection("pageContents")}
}

// Get single page
func (r PageRepository) ReadOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (page *models.PageModel, err error) {
	var _page models.PageModel

	if err = r.collection.FindOne(
		ctx, filter, opts...,
	).Decode(&_page); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_page, nil
}

// Get single page content
func (r PageContentRepository) ReadOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (
	pageContent *models.PageContentModel,
	err error,
) {
	var _pageContent models.PageContentModel

	if err = r.collection.FindOne(
		ctx, filter, opts...,
	).Decode(&_pageContent); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_pageContent, nil
}

// Get multiple pages
func (r PageRepository) ReadMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (pages []*models.PageModel, err error) {
	var (
		page   *models.PageModel
		cursor *mongo.Cursor
	)

	if cursor, err = r.collection.Find(
		ctx, filter, opts...,
	); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		page = &models.PageModel{}
		if err = cursor.Decode(page); err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}

	return pages, nil
}

// Count total pages
func (r PageRepository) Count(
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {

	return r.collection.CountDocuments(
		ctx, filter, opts...,
	)
}

// Save new page
func (r PageRepository) Save(
	ctx context.Context,
	page *models.PageModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID interface{}
		ok         bool
	)

	if insRes, err = r.collection.InsertOne(
		ctx, page,
	); err != nil {
		return err
	}
	if insertedID, ok = insRes.InsertedID.(primitive.ObjectID); !ok {
		return errors.New("unable to assert inserted uid")
	}
	if page.UID != insertedID {
		return errors.New("inserted uid is not same with database")
	}

	return nil
}

// Save new page content
func (r PageContentRepository) Save(
	ctx context.Context,
	pageContent *models.PageContentModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID interface{}
		ok         bool
	)

	if insRes, err = r.collection.InsertOne(
		ctx, pageContent,
	); err != nil {
		return err
	}
	if insertedID, ok = insRes.InsertedID.(primitive.ObjectID); !ok {
		return errors.New("unable to assert inserted uid")
	}
	if pageContent.UID != insertedID {
		return errors.New("inserted uid is not same with database")
	}

	return nil
}

// Update page
func (r PageRepository) Update(
	ctx context.Context,
	page *models.PageModel,
) (err error) {
	if _, err = r.collection.UpdateByID(
		ctx, page.UID, bson.M{"$set": page},
	); err != nil {
		return err
	}

	return nil
}

// Update page content
func (r PageContentRepository) Update(
	ctx context.Context,
	pageContent *models.PageContentModel,
) (err error) {
	if _, err = r.collection.UpdateByID(
		ctx, pageContent.UID, bson.M{"$set": pageContent},
	); err != nil {
		return err
	}

	return nil
}

// Delete page
func (r PageRepository) Delete(
	ctx context.Context,
	page *models.PageModel,
) (err error) {
	if _, err = r.collection.DeleteOne(
		ctx, bson.M{"_id": page.UID},
	); err != nil {
		return err
	}

	return nil
}

// Delete page content
func (r PageContentRepository) Delete(
	ctx context.Context,
	pageContent *models.PageContentModel,
) (err error) {
	if _, err = r.collection.DeleteOne(
		ctx, bson.M{"_id": pageContent.UID},
	); err != nil {
		return err
	}

	return err
}
