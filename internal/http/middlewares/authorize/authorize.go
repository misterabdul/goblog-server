package authorize

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
)

func Authorize(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
	level string,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me  *models.UserModel
			err error
		)

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
