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
	internalJwt "github.com/misterabdul/goblog-server/internal/pkg/jwt"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/pkg/hash"
)

func Refresh(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn      *mongo.Database
			userUid     primitive.ObjectID
			user        *models.UserModel
			refreshFlag bool
			err         error
		)
		userUid_s := c.GetString(authenticate.RefreshUserUid)
		tokenUid := c.GetString(authenticate.RefreshTokenUid)

		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if userUid, err = primitive.ObjectIDFromHex(userUid_s); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if user, err = repositories.GetUser(ctx, dbConn, bson.M{"_id": userUid}); err != nil {
			responses.Unauthenticated(c)
			return
		}

		newIssuedRefreshTokens := []models.IssuedToken{}
		for _, refreshToken := range user.IssuedRefreshTokens {
			if hash.Check(tokenUid, refreshToken.TokenUID) {
				refreshFlag = true
			} else {
				newIssuedRefreshTokens = append(newIssuedRefreshTokens, refreshToken)
			}
		}
		user.IssuedAccessTokens = newIssuedRefreshTokens
		if !refreshFlag {
			responses.Unauthenticated(c)
			return
		}

		accessTokenClaims, accessToken, err := internalJwt.IssueAccessToken(user)
		if err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		refreshTokenClaims, refreshToken, err := internalJwt.IssueRefreshToken(user)
		if err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		if err := saveToken(ctx, dbConn, user, accessTokenClaims, refreshTokenClaims); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responses.SignedIn(c, accessToken, accessTokenClaims, refreshToken, refreshTokenClaims)
	}
}
