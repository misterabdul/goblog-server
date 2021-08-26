package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	commentsCollectionName = "comments"
)

// Create the comments collection.
type CreateCommentsCollection struct {
}

func (m *CreateCommentsCollection) Name() string {
	return "03_create_comments_collections"
}

func (m *CreateCommentsCollection) Up(ctx context.Context, dbConn *mongo.Database) error {
	if err := dbConn.CreateCollection(ctx, commentsCollectionName); err != nil {
		return err
	}

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "postUid", Value: 1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "postSlug", Value: 1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "replies.email", Value: 1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "replies.name", Value: 1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "replies.createdAt", Value: 1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "replies.deletedAt", Value: 1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "createdAt", Value: 1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "deletedAt", Value: 1}},
			Options: nil,
		},
	}

	_, err := dbConn.Collection(commentsCollectionName).Indexes().CreateMany(ctx, indexes)

	return err
}

func (m *CreateCommentsCollection) Down(ctx context.Context, dbConn *mongo.Database) error {
	return dbConn.Collection(commentsCollectionName).Drop(ctx)
}
