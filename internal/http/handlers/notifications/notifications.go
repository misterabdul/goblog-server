package notifications

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/service"
)

func GetNotification(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel         = context.WithTimeout(context.Background(), maxCtxDuration)
			me                  *models.UserModel
			notification        *models.NotificationModel
			notificationId      primitive.ObjectID
			notificationIdQuery = c.Param("notification")
			err                 error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if notificationId, err = primitive.ObjectIDFromHex(notificationIdQuery); err != nil {
			responses.IncorrectNotificationId(c, err)
			return
		}
		if notification, err = svc.Notification.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"owner.username": me.Username},
				{"_id": notificationId}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if notification == nil {
			responses.NotFound(c, errors.New("notification not found"))
			return
		}

		responses.MyNotification(c, notification)
	}
}

func GetNotifications(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel   = context.WithTimeout(context.Background(), maxCtxDuration)
			me            *models.UserModel
			notifications []*models.NotificationModel
			err           error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if notifications, err = svc.Notification.GetMany(ctx, bson.M{
			"$and": []bson.M{
				{"owner.username": me.Username}}},
			internalGin.GetFindOptions(c),
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(notifications) == 0 {
			responses.NoContent(c)
			return
		}

		responses.MyNotifications(c, notifications)
	}
}

func ReadNotification(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel         = context.WithTimeout(context.Background(), maxCtxDuration)
			me                  *models.UserModel
			notification        *models.NotificationModel
			notificationId      primitive.ObjectID
			notificationIdQuery = c.Param("notification")
			err                 error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if notificationId, err = primitive.ObjectIDFromHex(notificationIdQuery); err != nil {
			responses.IncorrectNotificationId(c, err)
			return
		}
		if notification, err = svc.Notification.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"owner.username": me.Username},
				{"_id": notificationId}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if notification == nil {
			responses.NotFound(c, errors.New("notification not found"))
			return
		}
		if notification.ReadAt != nil {
			responses.NoContent(c)
			return
		}
		if err = svc.Notification.MarkOneRead(ctx, notification); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DeleteNotification(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel         = context.WithTimeout(context.Background(), maxCtxDuration)
			me                  *models.UserModel
			notification        *models.NotificationModel
			notificationId      primitive.ObjectID
			notificationIdQuery = c.Param("notification")
			err                 error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if notificationId, err = primitive.ObjectIDFromHex(notificationIdQuery); err != nil {
			responses.IncorrectNotificationId(c, err)
			return
		}
		if notification, err = svc.Notification.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"owner.username": me.Username},
				{"_id": notificationId}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if notification == nil {
			responses.NotFound(c, errors.New("notification not found"))
			return
		}
		if err = svc.Notification.DeleteOne(ctx, notification); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}
