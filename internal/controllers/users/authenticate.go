package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/repositories"
	refreshtokens "github.com/misterabdul/goblog-server/internal/repositories/refreshTokens"
	"github.com/misterabdul/goblog-server/internal/repositories/users"
)

func SignIn(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = users.GetUser(dbConn, "")

	c.JSON(http.StatusOK, gin.H{})
}

func SignInRefresh(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = refreshtokens.GetToken(dbConn, "")

	c.JSON(http.StatusOK, gin.H{})
}

func SignUp(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	user := users.CreateUser(dbConn, "")

	c.JSON(http.StatusOK, user)
}

func SignOut(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = refreshtokens.GetToken(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}
