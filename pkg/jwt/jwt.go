package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomClaims struct {
	KeyID     string `json:"kid,omitempty"`
	Id        string `json:"jti,omitempty"`
	Type      string `json:"typ,omitempty"`
	Subject   string `json:"sub,omitempty"`
	Audience  string `json:"aud,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
}

func (c CustomClaims) Valid() (err error) {
	var (
		vErr = new(jwt.ValidationError)
		now  = jwt.TimeFunc().Unix()
	)

	if c.ExpiresAt <= now {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("token is expired by %v", delta)
		vErr.Errors |= jwt.ValidationErrorExpired
	}
	if c.IssuedAt > now {
		vErr.Inner = fmt.Errorf("token used before issued")
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}
	if c.NotBefore > now {
		vErr.Inner = fmt.Errorf("token is not valid yet")
		vErr.Errors |= jwt.ValidationErrorNotValidYet
	}
	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}

func (c CustomClaims) GetExpiresAtSeconds() (expiresAtInSeconds int) {
	expiresAtTime := time.Unix(c.ExpiresAt, 0)
	diff := time.Until(expiresAtTime)

	return int(diff.Seconds())
}

func Issue(
	claimType string,
	subject string,
	duration time.Duration,
	secret string,
) (
	claims *CustomClaims,
	tokenString string,
	err error,
) {
	now := time.Now()
	tokenID := primitive.NewObjectID().Hex()
	claims = &CustomClaims{
		Id:        tokenID,
		Type:      claimType,
		Subject:   subject,
		Audience:  "goblog-client",
		Issuer:    "goblog-server",
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ExpiresAt: now.Add(duration).Unix()}
	if tokenString, err = IssueClaims(claims, secret); err != nil {
		return nil, "", err
	}

	return claims, tokenString, nil
}

func IssueClaims(claims *CustomClaims, secret string) (tokenString string, err error) {
	jwtClaim := jwt.NewWithClaims(jwt.SigningMethodHS512, *claims)
	tokenString, err = jwtClaim.SignedString([]byte(secret))

	return tokenString, err
}

func Check(tokenString string, secret string) (claims *CustomClaims, err error) {
	var (
		token     *jwt.Token
		rawClaims = CustomClaims{}
		ok        bool
	)

	if token, err = jwt.ParseWithClaims(
		tokenString,
		&rawClaims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		}); err != nil {
		return nil, err
	}
	if claims, ok = token.Claims.(*CustomClaims); !ok {
		return nil, errors.New("couldn't parse claims")
	}
	if err = claims.Valid(); err != nil {
		return nil, err
	}

	return claims, nil
}
