package repositories

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/database/models"
)

const (
	pageCollection        = "pages"
	pageContentCollection = "pageContents"
)

// Get single page
func ReadOnePage(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (page *models.PageModel, err error) {
	var (
		collection = dbConn.Collection(pageCollection)
		_page      models.PageModel
	)

	if err = collection.FindOne(ctx, filter, opts...).Decode(&_page); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_page, nil
}

// Get single page content
func ReadOnePageContent(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (pageContent *models.PageContentModel, err error) {
	var (
		collection   = dbConn.Collection(pageContentCollection)
		_pageContent models.PageContentModel
	)

	if err = collection.FindOne(ctx, filter, opts...).Decode(&_pageContent); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_pageContent, nil
}

// Get multiple pages
func ReadManyPages(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (pages []*models.PageModel, err error) {
	var (
		collection = dbConn.Collection(pageCollection)
		page       *models.PageModel
		cursor     *mongo.Cursor
	)

	if cursor, err = collection.Find(ctx, filter, opts...); err != nil {
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
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {
	var collection = dbConn.Collection(pageCollection)

	return collection.CountDocuments(
		ctx, filter, opts...)
}

// Save new page
func SaveOnePage(
	dbConn *mongo.Database,
	ctx context.Context,
	page *models.PageModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var (
		collection = dbConn.Collection(pageCollection)
		insRes     *mongo.InsertOneResult
		insertedID interface{}
		ok         bool
	)

	if insRes, err = collection.InsertOne(ctx, page, opts...); err != nil {
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
func SaveOnePageContent(
	dbConn *mongo.Database,
	ctx context.Context,
	pageContent *models.PageContentModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var (
		collection = dbConn.Collection(pageContentCollection)
		insRes     *mongo.InsertOneResult
		insertedID interface{}
		ok         bool
	)

	if insRes, err = collection.InsertOne(ctx, pageContent, opts...); err != nil {
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
func UpdateOnePage(
	dbConn *mongo.Database,
	ctx context.Context,
	page *models.PageModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var collection = dbConn.Collection(pageCollection)

	_, err = collection.UpdateOne(
		ctx, bson.M{"_id": page.UID}, bson.M{"$set": page}, opts...)

	return err
}

// Update page content
func UpdateOnePageContent(
	dbConn *mongo.Database,
	ctx context.Context,
	pageContent *models.PageContentModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var collection = dbConn.Collection(pageContentCollection)

	_, err = collection.UpdateOne(
		ctx, bson.M{"_id": pageContent.UID}, bson.M{"$set": pageContent}, opts...)

	return err
}

// Delete page
func DeleteOnePage(
	dbConn *mongo.Database,
	ctx context.Context,
	page *models.PageModel,
	opts ...*options.DeleteOptions,
) (err error) {
	var collection = dbConn.Collection(pageContentCollection)

	_, err = collection.DeleteOne(
		ctx, bson.M{"_id": page.UID}, opts...)

	return err
}

// Delete page content
func DeleteOnePageContent(
	dbConn *mongo.Database,
	ctx context.Context,
	pageContent *models.PageContentModel,
	opts ...*options.DeleteOptions,
) (err error) {
	var collection = dbConn.Collection(pageContentCollection)

	_, err = collection.DeleteOne(
		ctx, bson.M{"_id": pageContent.UID}, opts...)

	return err
}
