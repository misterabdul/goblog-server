package authentications

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	internalJwt "github.com/misterabdul/goblog-server/internal/pkg/jwt"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/pkg/hash"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

func SignIn(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			input              *forms.SignInForm
			user               *models.UserModel
			accessTokenClaims  *jwt.Claims
			refreshTokenClaims *jwt.Claims
			accessToken        string
			refreshToken       string
			err                error
		)

		if input, err = requests.GetSignInForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if user, err = repositories.GetUser(ctx, dbConn,
			bson.M{"$or": []bson.M{
				{"username": input.Username},
				{"email": input.Username},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if user == nil {
			responses.WrongSignIn(c, errors.New("incorrect username or email"))
			return
		}
		if !hash.Check(input.Password, user.Password) {
			responses.WrongSignIn(c, errors.New("incorrect password"))
			return
		}
		if accessTokenClaims, accessToken, err = internalJwt.IssueAccessToken(user); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if refreshTokenClaims, refreshToken, err = internalJwt.IssueRefreshToken(user); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = saveToken(ctx, dbConn, user, accessTokenClaims, refreshTokenClaims); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.SignedIn(c, accessToken, accessTokenClaims, refreshToken, refreshTokenClaims)
	}
}
