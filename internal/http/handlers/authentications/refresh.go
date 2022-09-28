package authentications

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	internalJwt "github.com/misterabdul/goblog-server/internal/pkg/jwt"
	"github.com/misterabdul/goblog-server/internal/service"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

// @Tags        Authentication
// @Summary     Refresh
// @Description Request new access token using refresh token.
// @Router      /v1/refresh [post]
// @Produce     application/json
// @Produce     application/msgpack
// @Header      200 {string} Set-Cookie
// @Success     200 {object} object{data=object{tokenType=string,accessToken=string}}
// @Success     201
// @Failure     401 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func Refresh(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel      = context.WithTimeout(context.Background(), maxCtxDuration)
			oldRefreshClaims *jwt.CustomClaims
			me               *models.UserModel
			newAccessClaims  *jwt.CustomClaims
			newAccessToken   string
			newRefreshClaims *jwt.CustomClaims
			newRefreshToken  string
			err              error
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
		if err = noteRevokeToken(ctx, svc, oldRefreshClaims, me); err != nil {
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
