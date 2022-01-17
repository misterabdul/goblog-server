package notifications

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/handlers/helpers"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

func GetNotification(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me                  *models.UserModel
			notification        *models.NotificationModel
			notificationId      primitive.ObjectID
			notificationIdQuery = c.Param("notification")
			err                 error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if notificationId, err = primitive.ObjectIDFromHex(notificationIdQuery); err != nil {
			responses.IncorrectNotificationId(c, err)
			return
		}
		if notification, err = repositories.GetNotification(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"owner.username": me.Username},
				{"_id": notificationId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if notification == nil {
			responses.NotFound(c, errors.New("notification not found"))
			return
		}

		responses.MyNotifiation(c, notification)
	}
}

func GetNotifications(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me            *models.UserModel
			notifications []*models.NotificationModel
			err           error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if notifications, err = repositories.GetNotifications(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"owner.username": me.Username}},
		}, helpers.GetFindOptions(c)); err != nil {
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
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me                  *models.UserModel
			notification        *models.NotificationModel
			notificationId      primitive.ObjectID
			notificationIdQuery = c.Param("notification")
			err                 error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if notificationId, err = primitive.ObjectIDFromHex(notificationIdQuery); err != nil {
			responses.IncorrectNotificationId(c, err)
			return
		}
		if notification, err = repositories.GetNotification(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"owner.username": me.Username},
				{"_id": notificationId}},
		}); err != nil {
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
		if err = repositories.ReadNotification(ctx, dbConn, notification); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DeleteNotification(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me                  *models.UserModel
			notification        *models.NotificationModel
			notificationId      primitive.ObjectID
			notificationIdQuery = c.Param("notification")
			err                 error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if notificationId, err = primitive.ObjectIDFromHex(notificationIdQuery); err != nil {
			responses.IncorrectNotificationId(c, err)
			return
		}
		if notification, err = repositories.GetNotification(ctx, dbConn, bson.M{
			"$and": []bson.M{
				{"owner.username": me.Username},
				{"_id": notificationId}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if notification == nil {
			responses.NotFound(c, errors.New("notification not found"))
			return
		}
		if err = repositories.DeleteNotification(ctx, dbConn, notification); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}
