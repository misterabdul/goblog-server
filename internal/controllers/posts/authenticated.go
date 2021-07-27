package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/internal/repositories/posts"
)

func GetMyPost(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	post := posts.GetPost(dbConn, "")

	c.JSON(http.StatusOK, post)
}

func GetMyPosts(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	posts := posts.GetPosts(dbConn, 10, "", "")

	c.JSON(http.StatusOK, posts)
}

func CreatePost(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = posts.CreatePost(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func PublishMyPost(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = posts.UpdatePost(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func DepublishMyPost(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = posts.UpdatePost(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func UpdateMyPost(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = posts.UpdatePost(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func TrashMyPost(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = posts.TrashPost(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func DeleteMyPost(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = posts.DeletePost(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func GetMyPostComment(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	comment := posts.GetComment(dbConn, "")

	c.JSON(http.StatusOK, comment)
}

func GetMyPostComments(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	comments := posts.GetComments(dbConn, 10, "", "")

	c.JSON(http.StatusOK, comments)
}

func TrashMyPostComment(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = posts.TrashComment(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func DeleteMyPostComment(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = posts.DeleteComment(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}
