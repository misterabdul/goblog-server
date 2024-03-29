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
	"github.com/misterabdul/goblog-server/internal/queue"
	"github.com/misterabdul/goblog-server/internal/queue/client"
	"github.com/misterabdul/goblog-server/internal/service"
)

// @Tags        Me
// @Summary     Update Me
// @Description Update my user data.
// @Router      /v1/auth/me [put]
// @Router      /v1/auth/me [patch]
// @Security    BearerAuth
// @Accept      application/json
// @Accept      application/msgpack
// @Produce     application/json
// @Produce     application/msgpack
// @Param       form body     object{firstname=string,lastname=string,username=string,email=string} true "Update me form"
// @Success     201
// @Failure     401  {object} object{message=string}
// @Failure     422  {object} object{message=string}
// @Failure     500  {object} object{message=string}
func UpdateMe(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			queueClient = client.GetClient()
			me          *models.UserModel
			updatedMe   *models.UserModel
			form        *forms.UpdateMeForm
			err         error
		)

		defer cancel()
		queueClient.Connect()
		defer queueClient.Disconnect()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if form, err = requests.GetUpdateMeForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if err = form.Validate(svc, ctx, me); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		updatedMe = form.ToUserModel(me)
		if err = svc.User.UpdateOne(ctx, updatedMe); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = queueClient.NewTask(queue.UpdateMe, updatedMe); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}
