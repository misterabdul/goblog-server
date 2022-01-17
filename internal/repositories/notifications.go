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

func getNotificationCollection(
	dbConn *mongo.Database,
) (notificationCollection *mongo.Collection) {
	return dbConn.Collection("notifications")
}

// Get single notification
func GetNotification(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.FindOneOptions,
) (notification *models.NotificationModel, err error) {
	var _notification models.NotificationModel

	if err = getNotificationCollection(dbConn).FindOne(
		ctx, filter, opts...,
	).Decode(&_notification); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_notification, nil
}

// Get multiple notifications
func GetNotifications(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.FindOptions,
) (notifications []*models.NotificationModel, err error) {
	var (
		notification *models.NotificationModel
		cursor       *mongo.Cursor
	)

	if cursor, err = getNotificationCollection(dbConn).Find(
		ctx, filter, opts...,
	); err != nil {
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
func SaveNotification(
	ctx context.Context,
	dbConn *mongo.Database,
	notification *models.NotificationModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID primitive.ObjectID
		ok         bool
	)

	if insRes, err = getNotificationCollection(dbConn).InsertOne(
		ctx, notification,
	); err != nil {
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
func UpdateNotification(
	ctx context.Context,
	dbConn *mongo.Database,
	notification *models.NotificationModel,
) (err error) {
	_, err = getNotificationCollection(dbConn).UpdateByID(
		ctx, notification.UID, bson.M{"$set": notification})

	return err
}

// Delete notification
func DeleteNotification(
	ctx context.Context,
	dbConn *mongo.Database,
	notification *models.NotificationModel,
) (err error) {
	_, err = getNotificationCollection(dbConn).DeleteOne(
		ctx, bson.M{"_id": notification.UID})

	return err
}
