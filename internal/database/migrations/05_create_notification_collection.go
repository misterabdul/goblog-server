package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	notificationCollectionName = "notification"
)

// Create the notification collection.
type CreateNotificationCollection struct {
}

func (m *CreateNotificationCollection) Name() string {
	return "05_create_notification_collections"
}

func (m *CreateNotificationCollection) Up(ctx context.Context, dbConn *mongo.Database) error {
	if err := dbConn.CreateCollection(ctx, notificationCollectionName); err != nil {
		return err
	}

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "owner.username", Value: 1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "createdAt", Value: -1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "deletedAt", Value: 1}},
			Options: nil,
		},
	}

	_, err := dbConn.Collection(notificationCollectionName).Indexes().CreateMany(ctx, indexes)

	return err
}

func (m *CreateNotificationCollection) Down(ctx context.Context, dbConn *mongo.Database) error {
	return dbConn.Collection(notificationCollectionName).Drop(ctx)
}
