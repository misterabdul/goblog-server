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

type CategoryRepository struct {
	collection *mongo.Collection
}

func NewCategoryRepository(
	dbConn *mongo.Database,
) *CategoryRepository {

	return &CategoryRepository{
		collection: dbConn.Collection("categories")}
}

// Get single category
func (r *CategoryRepository) ReadOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (category *models.CategoryModel, err error) {
	var _category models.CategoryModel

	if err = r.collection.FindOne(
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
func (r *CategoryRepository) ReadMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (categories []*models.CategoryModel, err error) {
	var (
		category *models.CategoryModel
		cursor   *mongo.Cursor
	)

	if cursor, err = r.collection.Find(
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

// Count total categories
func (r *CategoryRepository) Count(
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {

	return r.collection.CountDocuments(
		ctx, filter, opts...,
	)
}

// Save new category
func (r *CategoryRepository) Save(
	ctx context.Context,
	category *models.CategoryModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID primitive.ObjectID
		ok         bool
	)

	if insRes, err = r.collection.InsertOne(
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
func (r *CategoryRepository) Update(
	ctx context.Context,
	category *models.CategoryModel,
) (err error) {
	_, err = r.collection.UpdateByID(
		ctx, category.UID, bson.M{"$set": category})

	return err
}

// Delete category
func (r *CategoryRepository) Delete(
	ctx context.Context,
	category *models.CategoryModel,
) (err error) {
	_, err = r.collection.DeleteOne(
		ctx, bson.M{"_id": category.UID})

	return err
}
