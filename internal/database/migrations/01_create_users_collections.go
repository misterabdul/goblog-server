package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CreateUsersCollection struct {
}

func (m *CreateUsersCollection) Name() string {
	return "01_create_users_collections"
}

func (m *CreateUsersCollection) Up(ctx context.Context, dbConn *mongo.Database) error {
	if err := dbConn.CreateCollection(ctx, "users"); err != nil {
		return err
	}

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "createdAt", Value: -1}},
			Options: nil,
		},
	}

	_, err := dbConn.Collection("users").Indexes().CreateMany(ctx, indexes)

	return err
}

func (m *CreateUsersCollection) Down(ctx context.Context, dbConn *mongo.Database) error {
	return dbConn.Collection("users").Drop(ctx)
}
