package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/internal/repositories/posts"
)

func CreatePublicPostComment(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = posts.CreateComment(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func GetPublicPostComments(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	comments := posts.GetComments(dbConn, 10, "", "")

	c.JSON(http.StatusOK, comments)
}

func GetPublicPostComment(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	comment := posts.GetComment(dbConn, "")

	c.JSON(http.StatusOK, comment)
}
