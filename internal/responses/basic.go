package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/pkg/response"
)

func Basic(c *gin.Context, code int, obj interface{}) {
	acceptHeader := c.GetHeader("accept")
	switch {
	case acceptHeader == "application/json":
		c.JSON(code, obj)
		return
	case acceptHeader == "application/bson":
		response.BSON(c, code, obj)
		return
	case acceptHeader == "application/msgpack":
		response.MSGPACK(c, code, obj)
	}
	c.Data(http.StatusNotFound, "*/*", nil)
}
