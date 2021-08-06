package authenticate

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/pkg/jwt"
)

func AuthenticateRefresh() gin.HandlerFunc {

	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refreshToken")
		if err != nil {
			responses.Unauthenticated(c)
			c.Abort()
			return
		}

		claims, err := jwt.CheckRefreshToken(refreshToken)
		if err != nil {
			responses.Unauthenticated(c)
			c.Abort()
			return
		}

		c.Set(RefreshUserUid, claims.Payload.UserUID)
		c.Set(RefreshTokenUid, claims.TokenUID)

		c.Next()
	}
}
