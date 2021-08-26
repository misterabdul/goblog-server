package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	postCollectionName = "posts"
)

// Create the posts collection.
type CreatePostsCollection struct {
}

func (m *CreatePostsCollection) Name() string {
	return "02_create_posts_collections"
}

func (m *CreatePostsCollection) Up(ctx context.Context, dbConn *mongo.Database) error {
	if err := dbConn.CreateCollection(ctx, postCollectionName); err != nil {
		return err
	}

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "categories.slug", Value: 1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "tags", Value: 1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "author.username", Value: 1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "publishedAt", Value: -1}},
			Options: nil,
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
			Keys:    bson.D{{Key: "deletedAt", Value: -1}},
			Options: nil,
		},
	}

	_, err := dbConn.Collection(postCollectionName).Indexes().CreateMany(ctx, indexes)

	return err
}

func (m *CreatePostsCollection) Down(ctx context.Context, dbConn *mongo.Database) error {
	return dbConn.Collection(postCollectionName).Drop(ctx)
}
