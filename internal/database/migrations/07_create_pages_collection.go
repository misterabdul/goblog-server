package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const pageCollectionName = "pages"

// Create the page collection.
type CreatePagesCollection struct{}

func (m *CreatePagesCollection) Name() (collectionName string) {
	return "07_create_pages_collections"
}

func (m *CreatePagesCollection) Up(ctx context.Context, dbConn *mongo.Database) (err error) {
	if err = dbConn.CreateCollection(ctx, pageCollectionName); err != nil {
		return err
	}
	indexes := []mongo.IndexModel{{
		Keys: bson.D{
			{Key: "title", Value: "text"},
			{Key: "author.firstname", Value: "text"},
			{Key: "author.lastname", Value: "text"}},
		Options: nil,
	}, {
		Keys:    bson.D{{Key: "slug", Value: 1}},
		Options: options.Index().SetUnique(true),
	}, {
		Keys:    bson.D{{Key: "author.username", Value: 1}},
		Options: nil,
	}, {
		Keys:    bson.D{{Key: "publishedAt", Value: -1}},
		Options: nil,
	}, {
		Keys:    bson.D{{Key: "createdAt", Value: -1}},
		Options: nil,
	}, {
		Keys:    bson.D{{Key: "updatedAt", Value: -1}},
		Options: nil,
	}, {
		Keys:    bson.D{{Key: "deletedAt", Value: -1}},
		Options: nil,
	}}
	if _, err := dbConn.Collection(pageCollectionName).Indexes().
		CreateMany(ctx, indexes); err != nil {
		return err
	}

	return nil
}

func (m *CreatePagesCollection) Down(ctx context.Context, dbConn *mongo.Database) (err error) {
	return dbConn.Collection(pageCollectionName).Drop(ctx)
}
