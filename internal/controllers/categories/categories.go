package categories

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/internal/repositories/categories"
)

func GetPublicCategory(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	notification := categories.GetCategory(dbConn, "")

	c.JSON(http.StatusOK, notification)
}

func GetPublicCategorySlug(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	notification := categories.GetCategory(dbConn, "")

	c.JSON(http.StatusOK, notification)
}

func GetPublicCategories(c *gin.Context) {
	dbConn := repositories.GetDBConn("")
	notifications := categories.GetCategories(dbConn, 10, "createdAt", "desc")

	c.JSON(http.StatusOK, notifications)
}
