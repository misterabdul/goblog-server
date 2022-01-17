package service

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/models"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

// Get single notification
func (service *Service) GetNotification(
	filter interface{},
) (notification *models.NotificationModel, err error) {

	return repositories.GetNotification(
		service.ctx,
		service.dbConn,
		filter)
}

// Get multiple notifications
func (service *Service) GetNotifications(
	filter interface{},
) (notifications []*models.NotificationModel, err error) {

	return repositories.GetNotifications(
		service.ctx,
		service.dbConn,
		filter,
		internalGin.GetFindOptions(service.c))
}

// Create new notification
func (service *Service) CreateNotification(
	notification *models.NotificationModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	notification.UID = primitive.NewObjectID()
	notification.CreatedAt = now
	notification.DeletedAt = nil

	return repositories.SaveNotification(
		service.ctx,
		service.dbConn,
		notification)
}

// Mark the notification read
func (service *Service) ReadNotification(
	notification *models.NotificationModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	notification.ReadAt = now

	return repositories.UpdateNotification(
		service.ctx,
		service.dbConn,
		notification)
}

// Permanently delete notification
func (service *Service) DeleteNotification(
	notification *models.NotificationModel,
) (err error) {

	return repositories.DeleteNotification(
		service.ctx,
		service.dbConn,
		notification)
}
