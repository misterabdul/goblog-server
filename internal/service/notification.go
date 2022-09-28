package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/database/repositories"
)

type notification struct {
	dbConn *mongo.Database
}

func newNotificationService(
	dbConn *mongo.Database,
) (service *notification) {

	return &notification{dbConn: dbConn}
}

// Get single notification
func (s *notification) GetOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (notification *models.NotificationModel, err error) {

	return repositories.ReadOneNotification(
		s.dbConn, ctx, filter, opts...)
}

// Get multiple notifications
func (s *notification) GetMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (notifications []*models.NotificationModel, err error) {

	return repositories.ReadManyNotifications(
		s.dbConn, ctx, filter, opts...)
}

// Create new notification
func (s *notification) SaveOne(
	ctx context.Context,
	notification *models.NotificationModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	notification.UID = primitive.NewObjectID()
	notification.CreatedAt = now
	notification.DeletedAt = nil

	return repositories.SaveOneNotification(
		s.dbConn, ctx, notification, opts...)
}

// Mark the notification read
func (s *notification) MarkOneRead(
	ctx context.Context,
	notification *models.NotificationModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	notification.ReadAt = now

	return repositories.UpdateOneNotification(
		s.dbConn, ctx, notification, opts...)
}

// Permanently delete notification
func (s *notification) DeleteOne(
	ctx context.Context,
	notification *models.NotificationModel,
	opts ...*options.DeleteOptions,
) (err error) {

	return repositories.DeleteOneNotification(
		s.dbConn, ctx, notification, opts...)
}
