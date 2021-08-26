package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	categoryCollectionName = "category"
)

// Create the category collection.
type CreateCategoryCollection struct {
}

func (m *CreateCategoryCollection) Name() string {
	return "04_create_category_collections"
}

func (m *CreateCategoryCollection) Up(ctx context.Context, dbConn *mongo.Database) error {
	if err := dbConn.CreateCollection(ctx, categoryCollectionName); err != nil {
		return err
	}

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "createdAt", Value: -1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "updatedAt", Value: -1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "deletedAt", Value: 1}},
			Options: nil,
		},
	}

	_, err := dbConn.Collection(categoryCollectionName).Indexes().CreateMany(ctx, indexes)

	return err
}

func (m *CreateCategoryCollection) Down(ctx context.Context, dbConn *mongo.Database) error {
	return dbConn.Collection(categoryCollectionName).Drop(ctx)
}
