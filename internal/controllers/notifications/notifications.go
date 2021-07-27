package notifications

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/internal/repositories/notifications"
)

func GetNotification(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	notification := notifications.GetNotification(dbConn, "")

	c.JSON(http.StatusOK, notification)
}

func GetNotifications(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	notifications := notifications.GetNotifications(dbConn, 10, "createdAt", "desc")

	c.JSON(http.StatusOK, notifications)
}

func ReadNotification(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = notifications.CreateNotification(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}
