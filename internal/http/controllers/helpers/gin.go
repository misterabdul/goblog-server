package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetShowQuery(c *gin.Context) int {
	var (
		sQuery string
		query  int
		err    error
	)

	sQuery = c.DefaultQuery("show", "25")
	if query, err = strconv.Atoi(sQuery); err != nil {
		query = 25
	}

	return query
}

func GetOrderQuery(c *gin.Context) string {
	return c.DefaultQuery("order", "createdAt")
}

func GetAscQuery(c *gin.Context) bool {
	var (
		sQuery string
		query  bool
		err    error
	)

	sQuery = c.DefaultQuery("asc", "false")
	if query, err = strconv.ParseBool(sQuery); err != nil {
		query = false
	}

	return query
}
