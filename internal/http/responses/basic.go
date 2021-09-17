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
		return
	}
	c.Data(http.StatusNotFound, "*/*", nil)
}

func NoContent(c *gin.Context) {
	Basic(c, http.StatusNoContent, nil)
}

func BadRequest(c *gin.Context, msg string, err error) {
	Basic(c, http.StatusBadRequest, gin.H{"message": msg})
}

func NotFound(c *gin.Context, err error) {
	Basic(c, http.StatusNotFound, gin.H{"message": "not found"})
}

func FormIncorrect(c *gin.Context, err error) {
	Basic(c, http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
}

func InternalServerError(c *gin.Context, err error) {
	Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
}
