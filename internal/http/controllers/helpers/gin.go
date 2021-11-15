package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetFindOptions(c *gin.Context) (option *options.FindOptions) {
	var (
		show  = GetShowQuery(c)
		page  = GetPageQuery(c)
		order = GetOrderQuery(c)
		asc   = GetAscQuery(c)
	)

	*page = ((*page) - int64(1)) * (*show)

	return &options.FindOptions{
		Limit: show,
		Skip:  page,
		Sort:  bson.M{order: asc}}
}

func GetFindOptionsPost(c *gin.Context) (option *options.FindOptions) {
	var (
		show  = GetShowQuery(c)
		page  = GetPageQuery(c)
		order = GetPostOrderQuery(c)
		asc   = GetAscQuery(c)
	)

	*page = ((*page) - int64(1)) * (*show)

	return &options.FindOptions{
		Limit: show,
		Skip:  page,
		Sort:  bson.M{order: asc}}
}

func CreateFindOptions(
	show int,
	page int,
	order string,
	asc bool,
) (option *options.FindOptions) {
	var (
		show_i64 = int64(show)
		page_i64 = int64((page - 1) * show)
		asc_i    = 1
	)

	if !asc {
		asc_i = -1
	}

	return &options.FindOptions{
		Limit: &show_i64,
		Skip:  &page_i64,
		Sort:  bson.M{order: asc_i}}
}

func GetShowQuery(c *gin.Context) (show *int64) {
	var (
		sQuery string
		query  int64
		err    error
	)

	sQuery = c.DefaultQuery("show", "25")
	if query, err = strconv.ParseInt(sQuery, 10, 64); err != nil {
		query = 25
	}

	return &query
}

func GetPageQuery(c *gin.Context) (page *int64) {
	var (
		sQuery string
		query  int64
		err    error
	)

	sQuery = c.DefaultQuery("page", "1")
	if query, err = strconv.ParseInt(sQuery, 10, 64); err != nil {
		query = 25
	}

	return &query
}

func GetOrderQuery(c *gin.Context) (order string) {
	return c.DefaultQuery("order", "createdat")
}

func GetPostOrderQuery(c *gin.Context) (postOrder string) {
	return c.DefaultQuery("order", "publishedat")
}

func GetAscQuery(c *gin.Context) (asc int) {
	var (
		sQuery string
		query  bool
		err    error
	)

	sQuery = c.DefaultQuery("asc", "false")
	if query, err = strconv.ParseBool(sQuery); err != nil {
		query = false
	}
	if query {
		return 1
	}

	return -1
}
