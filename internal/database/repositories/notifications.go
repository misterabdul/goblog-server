package repositories

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/database/models"
)

const notificationCollection = "notifications"

// Get single notification
func ReadOneNotification(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (notification *models.NotificationModel, err error) {
	var (
		collection    = dbConn.Collection(notificationCollection)
		_notification models.NotificationModel
	)

	if err = collection.FindOne(ctx, filter, opts...).Decode(&_notification); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_notification, nil
}

// Get multiple notifications
func ReadManyNotifications(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (notifications []*models.NotificationModel, err error) {
	var (
		collection   = dbConn.Collection(notificationCollection)
		notification *models.NotificationModel
		cursor       *mongo.Cursor
	)

	if cursor, err = collection.Find(ctx, filter, opts...); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		notification = &models.NotificationModel{}
		if err = cursor.Decode(notification); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

// Save new notification
func SaveOneNotification(
	dbConn *mongo.Database,
	ctx context.Context,
	notification *models.NotificationModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var (
		collection = dbConn.Collection(notificationCollection)
		insRes     *mongo.InsertOneResult
		insertedID primitive.ObjectID
		ok         bool
	)

	if insRes, err = collection.InsertOne(ctx, notification, opts...); err != nil {
		return err
	}
	if insertedID, ok = insRes.InsertedID.(primitive.ObjectID); !ok {
		return errors.New("unable to assert inserted uid")
	}
	if notification.UID != insertedID {
		return errors.New("inserted uid is not same with database")
	}

	return nil
}

// Update notification
func UpdateOneNotification(
	dbConn *mongo.Database,
	ctx context.Context,
	notification *models.NotificationModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var collection = dbConn.Collection(notificationCollection)

	_, err = collection.UpdateOne(
		ctx, bson.M{"_id": notification.UID}, bson.M{"$set": notification}, opts...)

	return err
}

// Delete notification
func DeleteOneNotification(
	dbConn *mongo.Database,
	ctx context.Context,
	notification *models.NotificationModel,
	opts ...*options.DeleteOptions,
) (err error) {
	var collection = dbConn.Collection(notificationCollection)

	_, err = collection.DeleteOne(
		ctx, bson.M{"_id": notification.UID}, opts...)

	return err
}
