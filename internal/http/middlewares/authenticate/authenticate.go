package authenticate

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/pkg/jwt"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

const (
	AuthenticatedUserUid        = "AUTH_USER_UID"
	AuthenticatedTokenUid       = "AUTH_TOKEN_UID"
	AuthenticatedTokenExpiredAt = "AUTH_TOKEN_EXPIRED_AT"

	RefreshUserUid  = "REFRESH_USER_UID"
	RefreshTokenUid = "REFRESH_TOKEN_UID"
)

// Check the authentication status of given user.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		auth := c.GetHeader("Authorization")
		if !strings.Contains(auth, "Bearer ") {
			responses.Unauthenticated(c)
			c.Abort()
			return
		}
		auth = strings.ReplaceAll(auth, "Bearer ", "")

		claims, err := jwt.CheckAccessToken(auth)
		if err != nil {
			responses.Unauthenticated(c)
			c.Abort()
			return
		}

		var dbConn *mongo.Database
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		if isInvoked, err := checkForRevokedToken(ctx, dbConn, claims.Payload.UserUID, claims.TokenUID); isInvoked || err != nil {
			responses.Unauthenticated(c)
			c.Abort()
			return
		}

		c.Set(AuthenticatedUserUid, claims.Payload.UserUID)
		c.Set(AuthenticatedTokenUid, claims.TokenUID)
		c.Set(AuthenticatedTokenExpiredAt, claims.ExpiredAt)

		c.Next()
	}
}

func AuthenticateRefresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refreshToken")
		if err != nil {
			responses.Unauthenticated(c)
			c.Abort()
			return
		}

		claims, err := jwt.CheckRefreshToken(refreshToken)
		if err != nil {
			responses.Unauthenticated(c)
			c.Abort()
			return
		}

		c.Set(RefreshUserUid, claims.Payload.UserUID)
		c.Set(RefreshTokenUid, claims.TokenUID)

		c.Next()
	}
}

func checkForRevokedToken(ctx context.Context, dbConn *mongo.Database, userUid string, tokenUid string) (bool, error) {
	pUserUid, err := primitive.ObjectIDFromHex(userUid)
	if err != nil {
		return false, err
	}

	user, err := repositories.GetUser(ctx, dbConn, bson.M{"_id": pUserUid})
	if err != nil {
		return false, err
	}

	for _, revokedToken := range user.RevokedAccessTokens {
		if revokedToken.TokenUID == tokenUid {
			return true, nil
		}
	}

	return false, nil
}
