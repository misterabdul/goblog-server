package authentications

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

func SignOut(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			me     *models.UserModel
			claims *jwt.Claims
			err    error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}
		if claims, err = authenticate.GetAuthenticatedClaim(c); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "token claims not found"})
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		newIssuedAccessTokens := []models.IssuedToken{}
		for _, issuedAccessToken := range me.IssuedAccessTokens {
			if issuedAccessToken.TokenUID != claims.TokenUID {
				newIssuedAccessTokens = append(newIssuedAccessTokens, issuedAccessToken)
			}
		}
		me.IssuedAccessTokens = newIssuedAccessTokens
		me.RevokedAccessTokens = append(me.RevokedAccessTokens, models.RevokedToken{
			TokenUID: claims.TokenUID,
			Until:    primitive.NewDateTimeFromTime(claims.ExpiredAt),
		})
		if err := repositories.UpdateUser(ctx, dbConn, me); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responses.Basic(c, http.StatusNoContent, nil)
	}
}
