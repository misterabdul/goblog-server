package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const revokedTokenCollectionName = "revokedTokens"

// Create the revoke token collection.
type CreateRevokedTokenCollection struct{}

func (m *CreateRevokedTokenCollection) Name() (collectionName string) {
	return "06_create_revoke_token_collections"
}

func (m *CreateRevokedTokenCollection) Up(ctx context.Context, dbConn *mongo.Database) (err error) {
	if err = dbConn.CreateCollection(ctx, revokedTokenCollectionName); err != nil {
		return err
	}
	indexes := []mongo.IndexModel{{
		Keys:    bson.D{{Key: "expiresAt", Value: -1}},
		Options: nil,
	}, {
		Keys:    bson.D{{Key: "createdAt", Value: -1}},
		Options: nil,
	}, {
		Keys:    bson.D{{Key: "deletedAt", Value: 1}},
		Options: nil,
	}}
	if _, err := dbConn.Collection(revokedTokenCollectionName).Indexes().
		CreateMany(ctx, indexes); err != nil {
		return err
	}

	return nil
}

func (m *CreateRevokedTokenCollection) Down(ctx context.Context, dbConn *mongo.Database) (err error) {
	return dbConn.Collection(revokedTokenCollectionName).Drop(ctx)
}
