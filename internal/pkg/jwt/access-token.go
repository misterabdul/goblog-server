package jwt

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

const accessTokenTypeName = "access-token"

func IssueAccessToken(user *models.UserModel) (
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
	if duration_s, ok = os.LookupEnv("AUTH_ACCESS_DURATION"); !ok {
		duration_s = "60"
	}
	if duration, err = strconv.Atoi(duration_s); err != nil {
		duration = 60
	}
	if claims, tokenString, err = jwt.Issue(
		accessTokenTypeName,
		user.UID.Hex(),
		time.Duration(duration)*time.Minute,
		secret); err != nil {
		return nil, "", err
	}

	return claims, tokenString, nil
}

func CheckAccessToken(token string) (
	claims *jwt.CustomClaims,
	err error,
) {
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
	if claims.Type != accessTokenTypeName {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
