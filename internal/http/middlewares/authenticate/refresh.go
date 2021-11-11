package authenticate

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	internalJwt "github.com/misterabdul/goblog-server/internal/pkg/jwt"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

const (
	RefreshClaims = "REFRESH_CLAIMS"
	RefreshUser   = "REFRESH_USER"
)

func AuthenticateRefresh(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			refreshToken  string
			refreshClaims *jwt.CustomClaims
			revoked       bool
			userId        primitive.ObjectID
			me            *models.UserModel
			err           error
		)

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
		if revoked, err = checkRevokedToken(ctx, dbConn, refreshClaims); err != nil {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}
		if revoked {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}
		if userId, err = primitive.ObjectIDFromHex(refreshClaims.Subject); err != nil {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}
		if me, err = repositories.GetUser(ctx, dbConn, bson.M{"_id": userId}); err != nil {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}

		c.Set(RefreshClaims, *refreshClaims)
		c.Set(RefreshUser, *me)

		c.Next()
	}
}

func checkRevokedToken(
	ctx context.Context,
	dbConn *mongo.Database,
	refreshClaims *jwt.CustomClaims) (
	noted bool, err error) {
	var (
		tokenID          primitive.ObjectID
		revokedTokenData *models.RevokedTokenModel
	)

	if tokenID, err = primitive.ObjectIDFromHex(refreshClaims.Id); err != nil {
		return false, err
	}
	if revokedTokenData, err = repositories.GetRevokedToken(ctx, dbConn, bson.M{
		"_id": tokenID,
	}); err != nil {
		return false, err
	}
	if revokedTokenData != nil {
		return true, nil
	}

	return false, nil
}
