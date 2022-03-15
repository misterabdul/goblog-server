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
	Basic(c, http.StatusUnauthorized, gin.H{
		"message": "Unauthenticated."})
}

func WrongSignIn(c *gin.Context, err error) {
	Basic(c, http.StatusUnauthorized, gin.H{
		"message": "Wrong username or password."})
}

func SignedIn(
	c *gin.Context,
	accessToken string,
	accessClaims *jwt.CustomClaims,
	refreshToken string,
	refreshClaims *jwt.CustomClaims,
) {
	var (
		domain    string
		secured_s string
		secured   = false
		ok        bool
	)

	if domain, ok = os.LookupEnv("COOKIE_DOMAIN"); !ok {
		domain = ".localhost"
	}
	if secured_s, ok = os.LookupEnv("COOKIE_SECURE"); !ok {
		secured_s = "false"
	}
	if secured_s == "true" || secured_s == "TRUE" {
		secured = true
	}
	c.SetCookie(
		RefreshTokenCookieName,
		refreshToken,
		refreshClaims.GetExpiresAtSeconds(),
		"/api/v1/refresh",
		domain,
		secured,
		true)

	Basic(c, http.StatusOK, gin.H{
		"data": gin.H{
			"tokenType":   "Bearer",
			"accessToken": accessToken}})
}

func SigningOut(c *gin.Context) {
	var (
		domain    string
		secured_s string
		secured   = false
		ok        bool
	)

	if domain, ok = os.LookupEnv("COOKIE_DOMAIN"); !ok {
		domain = ".localhost"
	}
	if secured_s, ok = os.LookupEnv("COOKIE_SECURE"); !ok {
		secured_s = "false"
	}
	if secured_s == "true" || secured_s == "TRUE" {
		secured = true
	}
	c.SetCookie(
		RefreshTokenCookieName,
		"expired",
		0,
		"/api/v1/refresh",
		domain,
		secured,
		true)

	NoContent(c)
}
