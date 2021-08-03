package authentications

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SignUp(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{})
	}
}
