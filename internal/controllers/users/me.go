package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/internal/repositories/users"
)

func GetMe(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	user := users.GetUser(dbConn, "")

	c.JSON(http.StatusOK, user)
}

func UpdateMe(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	users.UpdateUser(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func UpdateMePassword(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	users.UpdateUser(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}
