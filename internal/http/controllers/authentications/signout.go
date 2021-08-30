package authentications

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

func SignOut(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me                    *models.UserModel
			claims                *jwt.Claims
			newIssuedAccessTokens = []models.IssuedToken{}
			err                   error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if claims, err = authenticate.GetAuthenticatedClaim(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
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
		if err = repositories.UpdateUser(ctx, dbConn, me); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}
