package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/models"
)

func MyNotification(c *gin.Context, notification *models.NotificationModel) {
	extractMyNotificationData(notification)
}

func MyNotifications(c *gin.Context, notifications []*models.NotificationModel) {
	var data []gin.H

	for _, notification := range notifications {
		data = append(data, extractMyNotificationData(notification))
	}
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func IncorrectNotificationId(c *gin.Context, err error) {
	Basic(c, http.StatusBadRequest, gin.H{
		"message": "incorrect notification id format"})
}

func extractMyNotificationData(notification *models.NotificationModel) (extracted gin.H) {
	return gin.H{
		"uid":       notification.UID.Hex(),
		"title":     notification.Title,
		"content":   notification.Content,
		"readAt":    notification.ReadAt,
		"createdAt": notification.CreatedAt}
}
