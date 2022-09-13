package authenticate

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

func GetAuthenticatedUser(c *gin.Context) (user *models.UserModel, err error) {
	var (
		_user   models.UserModel
		rawData interface{}
		ok      bool
	)

	if rawData, ok = c.Get(AuthenticatedUser); !ok {
		return nil, errors.New("no gin context data found for: " + AuthenticatedUser)
	}
	if _user, ok = rawData.(models.UserModel); !ok {
		return nil, errors.New("unable to assert user data")
	}

	return &_user, nil
}

func GetRefreshedUser(c *gin.Context) (user *models.UserModel, err error) {
	var (
		_user   models.UserModel
		rawData interface{}
		ok      bool
	)

	if rawData, ok = c.Get(RefreshUser); !ok {
		return nil, errors.New("no gin context data found for: " + AuthenticatedUser)
	}
	if _user, ok = rawData.(models.UserModel); !ok {
		return nil, errors.New("unable to assert user data")
	}

	return &_user, nil
}

func GetAuthenticatedClaim(c *gin.Context) (claims *jwt.CustomClaims, err error) {
	var (
		_claims jwt.CustomClaims
		rawData interface{}
		ok      bool
	)

	if rawData, ok = c.Get(AuthenticatedClaims); !ok {
		return nil, errors.New("no gin context data found for: " + AuthenticatedClaims)
	}
	if _claims, ok = rawData.(jwt.CustomClaims); !ok {
		return nil, errors.New("unable to assert jwt claims data")
	}

	return &_claims, nil
}

func GetRefreshClaims(c *gin.Context) (claims *jwt.CustomClaims, err error) {
	var (
		_claims jwt.CustomClaims
		rawData interface{}
		ok      bool
	)

	if rawData, ok = c.Get(RefreshClaims); !ok {
		return nil, errors.New("no gin context data found for: " + RefreshClaims)
	}
	if _claims, ok = rawData.(jwt.CustomClaims); !ok {
		return nil, errors.New("unable to assert jwt claims data")
	}

	return &_claims, nil
}
