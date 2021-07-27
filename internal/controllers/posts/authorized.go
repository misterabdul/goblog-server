package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/internal/repositories/posts"
)

func GetPost(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	post := posts.GetPost(dbConn, "")

	c.JSON(http.StatusOK, post)
}

func GetPosts(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	posts := posts.GetPosts(dbConn, 10, "", "")

	c.JSON(http.StatusOK, posts)
}

func PublishPost(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = posts.UpdatePost(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func DepublishPost(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = posts.UpdatePost(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func UpdatePost(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = posts.UpdatePost(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func TrashPost(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = posts.TrashPost(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func DeletePost(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = posts.DeletePost(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func GetPostComment(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	comment := posts.GetComment(dbConn, "")

	c.JSON(http.StatusOK, comment)
}

func GetPostComments(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	comments := posts.GetComments(dbConn, 10, "", "")

	c.JSON(http.StatusOK, comments)
}

func TrashPostComment(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = posts.TrashComment(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func DeletePostComment(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = posts.DeleteComment(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}
