package authenticate

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

func GetAuthenticatedUser(c *gin.Context) (*models.UserModel, error) {
	var (
		user    models.UserModel
		rawData interface{}
		ok      bool
	)
	if rawData, ok = c.Get(AuthenticatedUser); !ok {
		return nil, errors.New("no gin context data found for: " + AuthenticatedUser)
	}
	if user, ok = rawData.(models.UserModel); !ok {
		return nil, errors.New("unable to assert user data")
	}

	return &user, nil
}

func GetAuthenticatedClaim(c *gin.Context) (*jwt.Claims, error) {
	var (
		claims  jwt.Claims
		rawData interface{}
		ok      bool
	)
	if rawData, ok = c.Get(AuthenticatedClaims); !ok {
		return nil, errors.New("no gin context data found for: " + AuthenticatedClaims)
	}
	if claims, ok = rawData.(jwt.Claims); !ok {
		return nil, errors.New("unable to assert jwt claims data")
	}

	return &claims, nil
}
