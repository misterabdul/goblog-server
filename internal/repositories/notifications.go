package repositories

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/models"
)

func getNotificationCollection(dbConn *mongo.Database) *mongo.Collection {
	return dbConn.Collection("notifications")
}

// Get single notification
func GetNotification(ctx context.Context, dbConn *mongo.Database, filter interface{}) (*models.NotificationModel, error) {
	var notification models.NotificationModel
	if err := getNotificationCollection(dbConn).FindOne(ctx, filter).Decode(&notification); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &notification, nil
}

// Get multiple notifications
func GetNotifications(ctx context.Context, dbConn *mongo.Database, filter interface{}, show int, order string, asc bool) ([]*models.NotificationModel, error) {
	cursor, err := getNotificationCollection(dbConn).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []*models.NotificationModel
	for cursor.Next(ctx) {
		var notification models.NotificationModel
		if err := cursor.Decode(&notification); err != nil {
			return nil, err
		}
		notifications = append(notifications, &notification)
	}

	return notifications, nil
}

// Create new notification
func CreateNotification(ctx context.Context, dbConn *mongo.Database, notification *models.NotificationModel) error {
	now := primitive.NewDateTimeFromTime(time.Now())

	notification.UID = primitive.NewObjectID()
	notification.CreatedAt = now
	notification.DeletedAt = nil

	insRes, err := getNotificationCollection(dbConn).InsertOne(ctx, notification)
	if err != nil {
		return err
	}
	insertedID, ok := insRes.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("unable to assert inserted uid")
	}
	if notification.UID != insertedID {
		return errors.New("inserted uid is not same with database")
	}

	return nil
}

// Read notification
func ReadNotification(ctx context.Context, dbConn *mongo.Database, notification *models.NotificationModel) error {
	now := primitive.NewDateTimeFromTime(time.Now())

	notification.ReadAt = now

	_, err := getNotificationCollection(dbConn).UpdateByID(ctx, notification.UID, bson.M{"$set": notification})

	return err
}

// Permanently delete notification
func DeleteNotification(ctx context.Context, dbConn *mongo.Database, notification *models.NotificationModel) error {
	_, err := getNotificationCollection(dbConn).DeleteOne(ctx, bson.M{"_id": notification.UID})

	return err
}
