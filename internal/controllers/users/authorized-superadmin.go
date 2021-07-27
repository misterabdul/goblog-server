package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/internal/repositories/users"
)

func AdminizeUser(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = users.UpdateUser(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func DeadminizeUser(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = users.UpdateUser(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}
