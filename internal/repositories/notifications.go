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

type NotificationRepository struct {
	collection *mongo.Collection
}

func NewNotificationRepository(
	dbConn *mongo.Database,
) *NotificationRepository {

	return &NotificationRepository{
		collection: dbConn.Collection("notifications")}
}

// Get single notification
func (r NotificationRepository) ReadOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (notification *models.NotificationModel, err error) {
	var _notification models.NotificationModel

	if err = r.collection.FindOne(
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
func (r NotificationRepository) ReadMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (notifications []*models.NotificationModel, err error) {
	var (
		notification *models.NotificationModel
		cursor       *mongo.Cursor
	)

	if cursor, err = r.collection.Find(
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
func (r NotificationRepository) Save(
	ctx context.Context,
	notification *models.NotificationModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID primitive.ObjectID
		ok         bool
	)

	if insRes, err = r.collection.InsertOne(
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
func (r NotificationRepository) Update(
	ctx context.Context,
	notification *models.NotificationModel,
) (err error) {
	_, err = r.collection.UpdateByID(
		ctx, notification.UID, bson.M{"$set": notification})

	return err
}

// Delete notification
func (r NotificationRepository) Delete(
	ctx context.Context,
	notification *models.NotificationModel,
) (err error) {
	_, err = r.collection.DeleteOne(
		ctx, bson.M{"_id": notification.UID})

	return err
}
