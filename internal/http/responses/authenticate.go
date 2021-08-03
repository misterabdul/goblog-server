package responses

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/pkg/jwt"
)

func Unauthenticated(c *gin.Context) {
	Basic(c, http.StatusUnauthorized, gin.H{"message": "Unauthenticated."})
}

func SignedIn(c *gin.Context, accessToken string, accessTokenClaims *jwt.Claims, refreshToken string, refreshTokenClaims *jwt.Claims) {
	domain, ok := os.LookupEnv("COOKIE_DOMAIN")
	if !ok {
		domain = ".localhost"
	}

	secured_s, ok := os.LookupEnv("COOKIE_SECURE")
	if !ok {
		secured_s = "false"
	}
	secured := false
	if secured_s == "true" || secured_s == "TRUE" {
		secured = true
	}

	c.SetCookie(
		"refreshToken",
		refreshToken,
		refreshTokenClaims.ExpiredAt.Minute(),
		"/api/v1/refresh",
		domain,
		secured,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"tokenType":   "Bearer",
		"accessToken": accessToken,
	})
}
