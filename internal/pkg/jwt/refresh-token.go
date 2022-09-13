package jwt

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

const refreshTokenTypeName = "refresh-token"

func IssueRefreshToken(user *models.UserModel) (
	claims *jwt.CustomClaims,
	tokenString string,
	err error,
) {
	var (
		secret     string
		duration_s string
		duration   int
		ok         bool
	)

	if secret, ok = os.LookupEnv("AUTH_SECRET"); !ok {
		return nil, "", errors.New("unable to get authentication secret data")
	}
	if duration_s, ok = os.LookupEnv("AUTH_REFRESH_DURATION"); !ok {
		duration_s = "14"
	}
	if duration, err = strconv.Atoi(duration_s); err != nil {
		duration = 14
	}
	if claims, tokenString, err = jwt.Issue(
		refreshTokenTypeName,
		user.UID.Hex(),
		time.Duration(duration)*24*time.Hour,
		secret); err != nil {
		return nil, "", err
	}

	return claims, tokenString, nil
}

func CheckRefreshToken(token string) (claims *jwt.CustomClaims, err error) {
	var (
		secret string
		ok     bool
	)

	if secret, ok = os.LookupEnv("AUTH_SECRET"); !ok {
		return nil, errors.New("unable to get authentication secret data")
	}
	if claims, err = jwt.Check(token, secret); err != nil {
		return nil, err
	}
	if claims.Type != refreshTokenTypeName {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
