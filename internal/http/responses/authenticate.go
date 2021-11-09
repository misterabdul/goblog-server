package responses

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/pkg/jwt"
)

const (
	RefreshTokenCookieName = "refresh-token"
)

func Unauthenticated(c *gin.Context, err error) {
	Basic(c, http.StatusUnauthorized, gin.H{"message": "Unauthenticated."})
}

func WrongSignIn(c *gin.Context, err error) {
	Basic(c, http.StatusUnauthorized, gin.H{"message": "Wrong username or password."})
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
		RefreshTokenCookieName,
		refreshToken,
		refreshTokenClaims.ExpireDurationsInSeconds(),
		"/api/v1/refresh",
		domain,
		secured,
		true,
	)

	Basic(c, http.StatusOK, gin.H{
		"tokenType":   "Bearer",
		"accessToken": accessToken,
	})
}
