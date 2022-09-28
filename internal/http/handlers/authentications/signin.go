package authentications

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	internalJwt "github.com/misterabdul/goblog-server/internal/pkg/jwt"
	"github.com/misterabdul/goblog-server/internal/service"
	"github.com/misterabdul/goblog-server/pkg/hash"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

// @Tags        Authentication
// @Summary     Sign In
// @Description Do the signing in request.
// @Router      /v1/signin [post]
// @Accept      application/json
// @Accept      application/msgpack
// @Produce     application/json
// @Produce     application/msgpack
// @Param       form body     object{username=string,password=string} true "Login form"
// @Header      200  {string} Set-Cookie
// @Success     200  {object} object{data=object{tokenType=string,accessToken=string}}
// @Failure     401  {object} object{message=string}
// @Failure     422  {object} object{message=string}
// @Failure     500  {object} object{message=string}
func SignIn(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel   = context.WithTimeout(context.Background(), maxCtxDuration)
			input         *forms.SignInForm
			user          *models.UserModel
			accessClaims  *jwt.CustomClaims
			refreshClaims *jwt.CustomClaims
			accessToken   string
			refreshToken  string
			err           error
		)

		defer cancel()
		if input, err = requests.GetSignInForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if user, err = svc.User.GetOne(ctx, bson.M{
			"$or": []bson.M{
				{"username": bson.M{"$eq": input.Username}},
				{"email": bson.M{"$eq": input.Username}}},
		}); err != nil {
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
		if accessClaims, accessToken, err = internalJwt.IssueAccessToken(user); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if refreshClaims, refreshToken, err = internalJwt.IssueRefreshToken(user); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.SignedIn(
			c,
			accessToken,
			accessClaims,
			refreshToken,
			refreshClaims)
	}
}
