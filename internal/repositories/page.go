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

func getPageCollection(
	dbConn *mongo.Database,
) (pageCollection *mongo.Collection) {
	return dbConn.Collection("pages")
}

func getPageContentCollection(dbConn *mongo.Database,
) (pageContentCollection *mongo.Collection) {
	return dbConn.Collection("pageContents")
}

// Get single page
func GetPage(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.FindOneOptions,
) (page *models.PageModel, err error) {
	var _page models.PageModel

	if err = getPageCollection(dbConn).FindOne(
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
func GetPageContent(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.FindOneOptions,
) (
	pageContent *models.PageContentModel,
	err error,
) {
	var _pageContent models.PageContentModel

	if err = getPageContentCollection(dbConn).FindOne(
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
func GetPages(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.FindOptions,
) (pages []*models.PageModel, err error) {
	var (
		page   *models.PageModel
		cursor *mongo.Cursor
	)

	if cursor, err = getPageCollection(dbConn).Find(
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
func CountPages(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {

	return getPageCollection(dbConn).CountDocuments(
		ctx, filter, opts...,
	)
}

// Save new page
func SavePage(
	ctx context.Context,
	dbConn *mongo.Database,
	page *models.PageModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID interface{}
		ok         bool
	)

	if insRes, err = getPageCollection(dbConn).InsertOne(
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
func SavePageContent(
	ctx context.Context,
	dbConn *mongo.Database,
	pageContent *models.PageContentModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID interface{}
		ok         bool
	)

	if insRes, err = getPageContentCollection(dbConn).InsertOne(
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
func UpdatePage(
	ctx context.Context,
	dbConn *mongo.Database,
	page *models.PageModel,
) (err error) {
	if _, err = getPageCollection(dbConn).UpdateByID(
		ctx, page.UID, bson.M{"$set": page},
	); err != nil {
		return err
	}

	return nil
}

// Update page content
func UpdatePageContent(
	ctx context.Context,
	dbConn *mongo.Database,
	pageContent *models.PageContentModel,
) (err error) {
	if _, err = getPageContentCollection(dbConn).UpdateByID(
		ctx, pageContent.UID, bson.M{"$set": pageContent},
	); err != nil {
		return err
	}

	return nil
}

// Delete page
func DeletePage(
	ctx context.Context,
	dbConn *mongo.Database,
	page *models.PageModel,
) (err error) {
	if _, err = getPageCollection(dbConn).DeleteOne(
		ctx, bson.M{"_id": page.UID},
	); err != nil {
		return err
	}

	return nil
}

// Delete page content
func DeletePageContent(
	ctx context.Context,
	dbConn *mongo.Database,
	pageContent *models.PageContentModel,
) (err error) {
	if _, err = getPageContentCollection(dbConn).DeleteOne(
		ctx, bson.M{"_id": pageContent.UID},
	); err != nil {
		return err
	}

	return err
}
