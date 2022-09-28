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

const categoryCollection = "categories"

// Get single category
func ReadOneCategory(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (category *models.CategoryModel, err error) {
	var (
		collection = dbConn.Collection(categoryCollection)
		_category  models.CategoryModel
	)

	if err = collection.FindOne(ctx, filter, opts...).Decode(&_category); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_category, nil
}

// Get multiple categories
func ReadManyCategories(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (categories []*models.CategoryModel, err error) {
	var (
		collection = dbConn.Collection(categoryCollection)
		cursor     *mongo.Cursor
		category   *models.CategoryModel
	)

	if cursor, err = collection.Find(ctx, filter, opts...); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		category = &models.CategoryModel{}
		if err := cursor.Decode(category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// Count categories
func CountCategories(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {
	var collection = dbConn.Collection(categoryCollection)

	return collection.CountDocuments(
		ctx, filter, opts...)
}

// Save new category
func SaveOneCategory(
	dbConn *mongo.Database,
	ctx context.Context,
	category *models.CategoryModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var (
		collection = dbConn.Collection(categoryCollection)
		insRes     *mongo.InsertOneResult
		insertedID primitive.ObjectID
		ok         bool
	)

	if insRes, err = collection.InsertOne(ctx, category, opts...); err != nil {
		return err
	}
	if insertedID, ok = insRes.InsertedID.(primitive.ObjectID); !ok {
		return errors.New("unable to assert inserted uid")
	}
	if category.UID != insertedID {
		return errors.New("inserted uid is not same with database")
	}

	return nil
}

// Update category
func UpdateOneCategory(
	dbConn *mongo.Database,
	ctx context.Context,
	category *models.CategoryModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var collection = dbConn.Collection(categoryCollection)

	_, err = collection.UpdateOne(
		ctx, bson.M{"_id": category.UID}, bson.M{"$set": category}, opts...)

	return err
}

// Delete category
func DeleteOneCategory(
	dbConn *mongo.Database,
	ctx context.Context,
	category *models.CategoryModel,
	opts ...*options.DeleteOptions,
) (err error) {
	var collection = dbConn.Collection(categoryCollection)

	_, err = collection.DeleteOne(
		ctx, bson.M{"_id": category.UID}, opts...)

	return err
}
