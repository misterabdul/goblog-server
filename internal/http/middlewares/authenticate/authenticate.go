package authenticate

import (
	"context"
	"errors"
	"strings"
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
	AuthenticatedClaims = "AUTH_CLAIMS"
	AuthenticatedUser   = "AUTH_USER"
)

// Check the authentication status of given user.
func Authenticate(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			userService  = service.New(c, ctx, dbConn)
			me           *models.UserModel
			accessClaims *jwt.CustomClaims
			userId       primitive.ObjectID
			auth         string
			err          error
		)

		defer cancel()
		if auth = c.GetHeader("Authorization"); !strings.Contains(auth, "Bearer ") {
			responses.Unauthenticated(c, errors.New("no bearer type authorization header found"))
			c.Abort()
			return
		}
		auth = strings.ReplaceAll(auth, "Bearer ", "")
		if accessClaims, err = internalJwt.CheckAccessToken(auth); err != nil {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}
		if userId, err = primitive.ObjectIDFromHex(accessClaims.Subject); err != nil {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}
		if me, err = userService.GetUser(bson.M{"_id": userId}); err != nil {
			responses.Unauthenticated(c, err)
			c.Abort()
			return
		}
		c.Set(AuthenticatedClaims, *accessClaims)
		c.Set(AuthenticatedUser, *me)
		c.Next()
	}
}
