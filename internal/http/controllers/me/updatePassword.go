package me

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

func UpdateMePassword(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me        *models.UserModel
			updatedMe *models.UserModel
			form      *forms.UpdateMePasswordForm
			err       error
		)

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
		if err = repositories.UpdateUser(ctx, dbConn, updatedMe); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}
