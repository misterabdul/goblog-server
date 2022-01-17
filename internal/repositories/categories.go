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

func getCategoryCollection(
	dbConn *mongo.Database,
) (cateogryCollection *mongo.Collection) {
	return dbConn.Collection("categories")
}

// Get single category
func GetCategory(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.FindOneOptions,
) (category *models.CategoryModel, err error) {
	var _category models.CategoryModel

	if err = getCategoryCollection(dbConn).FindOne(
		ctx, filter, opts...,
	).Decode(&_category); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_category, nil
}

// Get multiple categories
func GetCategories(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.FindOptions,
) (categories []*models.CategoryModel, err error) {
	var (
		category *models.CategoryModel
		cursor   *mongo.Cursor
	)

	if cursor, err = getCategoryCollection(dbConn).Find(
		ctx, filter, opts...,
	); err != nil {
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

// Save new category
func SaveCategory(
	ctx context.Context,
	dbConn *mongo.Database,
	category *models.CategoryModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID primitive.ObjectID
		ok         bool
	)

	if insRes, err = getCategoryCollection(dbConn).InsertOne(
		ctx, category,
	); err != nil {
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
func UpdateCategory(
	ctx context.Context,
	dbConn *mongo.Database,
	category *models.CategoryModel,
) (err error) {
	_, err = getCategoryCollection(dbConn).UpdateByID(
		ctx, category.UID, bson.M{"$set": category})

	return err
}

// Delete category
func DeleteCategory(
	ctx context.Context,
	dbConn *mongo.Database,
	category *models.CategoryModel,
) (err error) {
	_, err = getCategoryCollection(dbConn).DeleteOne(
		ctx, bson.M{"_id": category.UID})

	return err
}
