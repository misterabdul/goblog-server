package users

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
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	internalJwt "github.com/misterabdul/goblog-server/internal/pkg/jwt"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/pkg/hash"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

func SignIn(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			input  *models.SignInModel
			dbConn *mongo.Database
			user   *models.UserModel
			err    error
		)

		if input, err = requests.GetSignInModel(c); err != nil {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if user, err = repositories.GetUser(ctx, dbConn, bson.M{"$or": []bson.M{
			{"username": input.Username},
			{"email": input.Username}}}); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "Wrong username or password."})
			return
		}
		if !hash.Check(input.Password, user.Password) {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "Wrong username or password."})
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

func SignInRefresh(maxCtxDuration time.Duration) gin.HandlerFunc {

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

func SignUp(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{})
	}
}

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

func saveToken(ctx context.Context, dbConn *mongo.Database, user *models.UserModel, accessToken *jwt.Claims, refreshToken *jwt.Claims) error {
	hashedRefreshToken, err := hash.Make(refreshToken.TokenUID)
	if err != nil {
		return err
	}

	user.IssuedRefreshTokens = append(user.IssuedRefreshTokens, models.IssuedToken{
		TokenUID:  hashedRefreshToken,
		Client:    "",
		IssuedAt:  primitive.NewDateTimeFromTime(refreshToken.IssuedAt),
		ExpiredAt: primitive.NewDateTimeFromTime(refreshToken.ExpiredAt),
	})
	user.IssuedAccessTokens = append(user.IssuedAccessTokens, models.IssuedToken{
		TokenUID:  accessToken.TokenUID,
		Client:    "",
		IssuedAt:  primitive.NewDateTimeFromTime(accessToken.IssuedAt),
		ExpiredAt: primitive.NewDateTimeFromTime(accessToken.ExpiredAt),
	})

	return repositories.UpdateUser(ctx, dbConn, user)
}
