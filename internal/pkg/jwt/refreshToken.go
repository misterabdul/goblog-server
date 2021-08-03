package jwt

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/pkg/crypto"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

func IssueRefreshToken(user *models.UserModel) (*jwt.Claims, string, error) {
	secret, ok := os.LookupEnv("AUTH_SECRET")
	if !ok {
		return nil, "", errors.New("unable to get authentication secret data")
	}

	duration_s, ok := os.LookupEnv("AUTH_REFRESH_DURATION")
	if !ok {
		duration_s = "14"
	}

	duration, err := strconv.Atoi(duration_s)
	if err != nil {
		duration = 14
	}

	payload := jwt.Payload{
		UserUID:  user.UID.Hex(),
		Username: user.Username,
	}

	tokenID, err := crypto.GenerateRandomString(64)
	if err != nil {
		return nil, "", err
	}

	claims := jwt.Claims{
		Payload:   payload,
		TokenUID:  tokenID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Duration(duration) * (24 * time.Hour)),
	}
	var token string

	if _, token, err = jwt.IssueClaims(claims, secret); err != nil {
		return nil, "", err
	}

	return &claims, token, nil
}

func CheckRefreshToken(token string) (claims jwt.Claims, err error) {
	secret, ok := os.LookupEnv("AUTH_SECRET")
	if !ok {
		return jwt.Claims{}, errors.New("unable to get authentication secret data")
	}

	if claims, err = jwt.Check(token, secret); err != nil {
		return jwt.Claims{}, err
	}

	return claims, nil
}
