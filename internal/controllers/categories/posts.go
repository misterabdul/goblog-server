package categories

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/internal/repositories/posts"
)

func GetPublicCategoryPosts(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	posts := posts.GetPosts(dbConn, 10, "", "")

	c.JSON(http.StatusOK, posts)
}

func GetPublicCategoryPost(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	post := posts.GetPost(dbConn, "")

	c.JSON(http.StatusOK, post)
}
