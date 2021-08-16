package authorize

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
)

func Authorize(maxCtxDuration time.Duration, level string) gin.HandlerFunc {

	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me  *models.UserModel
			err error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			c.Abort()
			return
		}
		if !CheckRoles(me, GetRole(level)) {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "unauthorized action"})
			c.Abort()
			return
		}

		c.Next()
	}
}
