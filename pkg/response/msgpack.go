package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack/v5"
)

// MsgPack http response for gin.
func MSGPACK(c *gin.Context, code int, obj interface{}) {
	var (
		msgpackData []byte
		err         error
	)

	if msgpackData, err = msgpack.Marshal(obj); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error()})
	}
	c.Data(code, "application/msgpack", msgpackData)
}
