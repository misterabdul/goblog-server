package authorize

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/service"
)

func Authorize(
	maxCtxDuration time.Duration,
	svc *service.Service,
	level string,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			_, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			me        *models.UserModel
			err       error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.UnauthorizedAction(c, err)
			c.Abort()
			return
		}
		if !CheckRoles(me, GetRole(level)) {
			responses.UnauthorizedAction(c, errors.New("unauthorized action"))
			c.Abort()
			return
		}
		c.Next()
	}
}
