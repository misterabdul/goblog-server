package me

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/service"
)

// @Tags        Me
// @Summary     Update Me Password
// @Description Update my password.
// @Router      /v1/auth/me/password [put]
// @Router      /v1/auth/me/password [patch]
// @Security    BearerAuth
// @Accept      application/json
// @Accept      application/msgpack
// @Produce     application/json
// @Produce     application/msgpack
// @Param       form body object{newpassword=string,newpasswordconfirm=string} true "Update me password form"
// @Success     201
// @Failure     401 {object} object{message=string}
// @Failure     422 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func UpdateMePassword(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			me          *models.UserModel
			updatedMe   *models.UserModel
			form        *forms.UpdateMePasswordForm
			err         error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if form, err = requests.GetUpdateMePasswordForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if err = form.Validate(me); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if updatedMe, err = form.ToUserModel(me); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = svc.User.UpdateOne(ctx, updatedMe); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}
