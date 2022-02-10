package authenticate

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	internalJwt "github.com/misterabdul/goblog-server/internal/pkg/jwt"
	"github.com/misterabdul/goblog-server/internal/service"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

const (
	RefreshClaims = "REFRESH_CLAIMS"
	RefreshUser   = "REFRESH_USER"
)

func AuthenticateRefresh(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel   = context.WithTimeout(context.Background(), maxCtxDuration)
			userService   = service.New(c, ctx, dbConn)
			refreshToken  string
			refreshClaims *jwt.CustomClaims
			revoked       bool
			userUid       primitive.ObjectID
			me            *models.UserModel
			err           error
		)

		defer cancel()
		if refreshToken, err = c.Cookie(responses.RefreshTokenCookieName); err != nil {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}
		if refreshClaims, err = internalJwt.CheckRefreshToken(refreshToken); err != nil {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}
		if revoked, err = checkRevokedToken(userService, refreshClaims); err != nil {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}
		if revoked {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}
		if userUid, err = primitive.ObjectIDFromHex(refreshClaims.Subject); err != nil {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}
		if me, err = userService.GetUser(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": userUid}}},
		}); err != nil {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}
		if me == nil {
			responses.Unauthenticated(c, errors.New("user not found"))
			c.Abort()
			return
		}
		c.Set(RefreshClaims, *refreshClaims)
		c.Set(RefreshUser, *me)
		c.Next()
	}
}

func checkRevokedToken(
	userService *service.Service,
	refreshClaims *jwt.CustomClaims,
) (noted bool, err error) {
	var (
		tokenUid         primitive.ObjectID
		revokedTokenData *models.RevokedTokenModel
	)

	if tokenUid, err = primitive.ObjectIDFromHex(refreshClaims.Id); err != nil {
		return false, err
	}
	if revokedTokenData, err = userService.GetRevokedToken(bson.M{
		"_id": tokenUid,
	}); err != nil {
		return false, err
	}
	if revokedTokenData != nil {
		return true, nil
	}

	return false, nil
}
