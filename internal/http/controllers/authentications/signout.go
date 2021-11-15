package authentications

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

func SignOut(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me            *models.UserModel
			refreshClaims *jwt.CustomClaims
			err           error
		)

		if me, err = authenticate.GetRefreshedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if refreshClaims, err = authenticate.GetRefreshClaims(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if err = noteRevokeToken(ctx, dbConn, refreshClaims, me); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}
