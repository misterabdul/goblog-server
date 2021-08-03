package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payload struct {
	UserUID  string `json:"userUid"`
	Username string `json:"username"`
}

type Claims struct {
	Payload
	TokenUID  string    `json:"tokenUid"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

func (c Claims) Valid() error {
	now := time.Now()

	if !c.IssuedAt.Before(now) {
		return errors.New("token is invalid")
	}
	if !c.ExpiredAt.After(now) {
		return errors.New("token is invalid")
	}

	return nil
}

func Issue(payload Payload, duration time.Duration, secret string) (tokenID string, token string, err error) {
	tokenID = primitive.NewObjectID().Hex()
	claims := Claims{
		payload,
		tokenID,
		time.Now(),
		time.Now().Add(duration),
	}
	return IssueClaims(claims, secret)
}

func IssueClaims(claims Claims, secret string) (tokenID string, token string, err error) {
	jwtClaim := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err = jwtClaim.SignedString([]byte(secret))

	return tokenID, token, err
}

func Check(tokenString string, secret string) (claims Claims, err error) {
	claims = Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return Claims{}, err
	}

	_claims, ok := token.Claims.(*Claims)
	if !ok {
		return Claims{}, errors.New("couldn't parse claims")
	}

	if err := _claims.Valid(); err != nil {
		return Claims{}, err
	}

	return *_claims, nil
}
