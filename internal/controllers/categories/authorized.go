package categories

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/internal/repositories/categories"
)

func GetCategories(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	categories := categories.GetCategories(dbConn, 10, "", "")

	c.JSON(http.StatusOK, categories)
}

func GetCategory(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	category := categories.GetCategory(dbConn, "")

	c.JSON(http.StatusOK, category)
}

func CreateCategory(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = categories.CreateCategory(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func UpdateCategory(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = categories.UpdateCategory(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func TrashCategory(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = categories.TrashCategory(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}

func DeleteCategory(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	_ = categories.DeleteCategory(dbConn, "")

	c.JSON(http.StatusNoContent, gin.H{})
}
