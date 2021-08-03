package authentications

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

func SignOut(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			user   *models.UserModel
			err    error
		)

		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		userUid, err := primitive.ObjectIDFromHex(c.GetString(authenticate.AuthenticatedUserUid))
		if err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}

		if user, err = repositories.GetUser(ctx, dbConn, bson.M{"_id": userUid}); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}

		tokenUid := c.GetString(authenticate.AuthenticatedTokenUid)
		newIssuedAccessTokens := []models.IssuedToken{}
		for _, issuedAccessToken := range user.IssuedAccessTokens {
			if issuedAccessToken.TokenUID != tokenUid {
				newIssuedAccessTokens = append(newIssuedAccessTokens, issuedAccessToken)
			}
		}
		user.IssuedAccessTokens = newIssuedAccessTokens

		tokenExpiredAt := c.GetTime(authenticate.AuthenticatedTokenExpiredAt)
		user.RevokedAccessTokens = append(user.RevokedAccessTokens, models.RevokedToken{
			TokenUID: tokenUid,
			Until:    primitive.NewDateTimeFromTime(tokenExpiredAt),
		})

		if err := repositories.UpdateUser(ctx, dbConn, user); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}

		responses.Basic(c, http.StatusNoContent, nil)
	}
}
