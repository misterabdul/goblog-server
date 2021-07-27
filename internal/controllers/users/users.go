package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/internal/repositories/users"
)

func GetPublicUser(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	user := users.GetUser(dbConn, "")

	c.JSON(http.StatusOK, user)
}

func GetPublicUsers(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	users := users.GetUsers(dbConn, 10, "createdAt", "desc")

	c.JSON(http.StatusOK, users)
}
