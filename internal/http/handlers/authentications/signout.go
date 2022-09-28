package authentications

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/service"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

// @Tags        Authentication
// @Summary     Sign Out
// @Description Do the signing out request to revoke access token & refresh token.
// @Router      /v1/signout [post]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Header      200 {string} Set-Cookie
// @Success     201
// @Failure     401 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func SignOut(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel   = context.WithTimeout(context.Background(), maxCtxDuration)
			me            *models.UserModel
			refreshClaims *jwt.CustomClaims
			err           error
		)

		defer cancel()
		if me, err = authenticate.GetRefreshedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if refreshClaims, err = authenticate.GetRefreshClaims(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if err = noteRevokeToken(ctx, svc, refreshClaims, me); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.SigningOut(c)
	}
}
