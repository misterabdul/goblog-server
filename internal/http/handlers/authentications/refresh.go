package authentications

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	internalJwt "github.com/misterabdul/goblog-server/internal/pkg/jwt"
	"github.com/misterabdul/goblog-server/internal/service"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

func Refresh(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel         = context.WithTimeout(context.Background(), maxCtxDuration)
			revokedTokenService = service.NewRevokedTokenService(c, ctx, dbConn)
			oldRefreshClaims    *jwt.CustomClaims
			me                  *models.UserModel
			newAccessClaims     *jwt.CustomClaims
			newAccessToken      string
			newRefreshClaims    *jwt.CustomClaims
			newRefreshToken     string
			err                 error
		)

		defer cancel()
		if oldRefreshClaims, err = authenticate.GetRefreshClaims(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if me, err = authenticate.GetRefreshedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if err = noteRevokeToken(revokedTokenService, oldRefreshClaims, me); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if newAccessClaims, newAccessToken, err = internalJwt.
			IssueAccessToken(me); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if newRefreshClaims, newRefreshToken, err = internalJwt.
			IssueRefreshToken(me); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.SignedIn(
			c,
			newAccessToken,
			newAccessClaims,
			newRefreshToken,
			newRefreshClaims)
	}
}
