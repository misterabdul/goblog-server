package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/internal/repositories/users"
)

func GetUsers(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	users := users.GetUsers(dbConn, 10, "", "")

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	user := users.GetUser(dbConn, "")

	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	users.CreateUser(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func UpdateUser(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	users.UpdateUser(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func TrashUser(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	users.TrashUser(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func DeleteUser(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	users.DeleteUser(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}
