package authenticate

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	internalJwt "github.com/misterabdul/goblog-server/internal/pkg/jwt"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

const (
	AuthenticatedClaims = "AUTH_CLAIMS"
	AuthenticatedUser   = "AUTH_USER"

	RefreshUserUid  = "REFRESH_USER_UID"
	RefreshTokenUid = "REFRESH_TOKEN_UID"
)

// Check the authentication status of given user.
func Authenticate(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			me     *models.UserModel
			claims jwt.Claims
			userId primitive.ObjectID
			auth   string
			err    error
		)

		if auth = c.GetHeader("Authorization"); !strings.Contains(auth, "Bearer ") {
			responses.Unauthenticated(c, errors.New("no bearer type authorization header found"))
			c.Abort()
			return
		}
		auth = strings.ReplaceAll(auth, "Bearer ", "")

		if claims, err = internalJwt.CheckAccessToken(auth); err != nil {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}
		if userId, err = primitive.ObjectIDFromHex(claims.Payload.UserUID); err != nil {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if me, err = repositories.GetUser(ctx, dbConn, bson.M{"_id": userId}); err != nil {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}
		for _, revokedToken := range me.RevokedAccessTokens {
			if revokedToken.TokenUID == claims.TokenUID {
				responses.Unauthenticated(c, errors.New("token already revoked"))
				c.Abort()
				return
			}
		}
		c.Set(AuthenticatedClaims, claims)
		c.Set(AuthenticatedUser, *me)

		c.Next()
	}
}
