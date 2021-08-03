package jwt

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/pkg/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IssueAccessToken(user *models.UserModel) (*jwt.Claims, string, error) {
	secret, ok := os.LookupEnv("AUTH_SECRET")
	if !ok {
		return nil, "", errors.New("unable to get authentication secret data")
	}

	duration_s, ok := os.LookupEnv("AUTH_ACCESS_DURATION")
	if !ok {
		duration_s = "60"
	}

	duration, err := strconv.Atoi(duration_s)
	if err != nil {
		duration = 60
	}

	payload := jwt.Payload{
		UserUID:  user.UID.Hex(),
		Username: user.Username,
	}

	claims := jwt.Claims{
		Payload:   payload,
		TokenUID:  primitive.NewObjectID().Hex(),
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Duration(duration) * time.Minute),
	}
	var token string

	if _, token, err = jwt.IssueClaims(claims, secret); err != nil {
		return nil, "", err
	}

	return &claims, token, nil
}

func CheckAccessToken(token string) (claims jwt.Claims, err error) {
	secret, ok := os.LookupEnv("AUTH_SECRET")
	if !ok {
		return jwt.Claims{}, errors.New("unable to get authentication secret data")
	}

	if claims, err = jwt.Check(token, secret); err != nil {
		return jwt.Claims{}, err
	}

	return claims, nil
}
