package service

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/database/repositories"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
)

type NotificationService struct {
	c          *gin.Context
	ctx        context.Context
	dbConn     *mongo.Database
	repository *repositories.NotificationRepository
}

func NewNotificationService(
	c *gin.Context,
	ctx context.Context,
	dbConn *mongo.Database,
) *NotificationService {

	return &NotificationService{
		c:          c,
		ctx:        ctx,
		dbConn:     dbConn,
		repository: repositories.NewNotificationRepository(dbConn)}
}

// Get single notification
func (s *NotificationService) GetNotification(
	filter interface{},
) (notification *models.NotificationModel, err error) {

	return s.repository.ReadOne(
		s.ctx, filter)
}

// Get multiple notifications
func (s *NotificationService) GetNotifications(
	filter interface{},
) (notifications []*models.NotificationModel, err error) {

	return s.repository.ReadMany(
		s.ctx, filter,
		internalGin.GetFindOptions(s.c))
}

// Create new notification
func (s *NotificationService) CreateNotification(
	notification *models.NotificationModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	notification.UID = primitive.NewObjectID()
	notification.CreatedAt = now
	notification.DeletedAt = nil

	return s.repository.Save(
		s.ctx, notification)
}

// Mark the notification read
func (s *NotificationService) ReadNotification(
	notification *models.NotificationModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	notification.ReadAt = now

	return s.repository.Update(
		s.ctx, notification)
}

// Permanently delete notification
func (s *NotificationService) DeleteNotification(
	notification *models.NotificationModel,
) (err error) {

	return s.repository.Delete(
		s.ctx, notification)
}
