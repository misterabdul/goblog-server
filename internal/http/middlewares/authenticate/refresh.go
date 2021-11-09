package authenticate

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/responses"
	internalJwt "github.com/misterabdul/goblog-server/internal/pkg/jwt"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

func AuthenticateRefresh() gin.HandlerFunc {

	return func(c *gin.Context) {
		var (
			refreshToken string
			claims       jwt.Claims
			err          error
		)

		if refreshToken, err = c.Cookie(responses.RefreshTokenCookieName); err != nil {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}

		if claims, err = internalJwt.CheckRefreshToken(refreshToken); err != nil {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}

		c.Set(RefreshUserUid, claims.Payload.UserUID)
		c.Set(RefreshTokenUid, claims.TokenUID)

		c.Next()
	}
}
