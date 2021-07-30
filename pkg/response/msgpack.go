package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack/v5"
)

func MSGPACK(c *gin.Context, code int, obj interface{}) {
	msgpackData, err := msgpack.Marshal(obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.Data(code, "application/msgpack", msgpackData)
}
