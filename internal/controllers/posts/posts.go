package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/internal/repositories/posts"
)

func GetPublicPost(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	post := posts.GetPost(dbConn, "")

	c.JSON(http.StatusOK, post)
}

func GetPublicPostSlug(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	post := posts.GetPost(dbConn, "")

	c.JSON(http.StatusOK, post)
}

func GetPublicPosts(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	posts := posts.GetPosts(dbConn, 10, "createdAt", "desc")

	c.JSON(http.StatusOK, posts)
}
