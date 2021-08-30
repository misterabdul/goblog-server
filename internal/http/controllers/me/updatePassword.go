package me

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/pkg/hash"
)

func UpdateMePassword(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me          *models.UserModel
			form        *forms.UpdateMePasswordForm
			newPassword string
			err         error
		)

		if form, err = requests.GetUpdateMePasswordForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if form.NewPassword != form.NewPasswordConfirm {
			responses.FormIncorrect(c, errors.New("new password confirmation not match"))
			return
		}
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if hash.Check(newPassword, me.Password) {
			responses.FormIncorrect(c, errors.New("your password is same as old one"))
			return
		}
		if newPassword, err = hash.Make(form.NewPassword); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		me.Password = newPassword
		if err = repositories.UpdateUser(ctx, dbConn, me); err != nil {
			var writeErr mongo.WriteException
			if errors.As(err, &writeErr) {
				responses.FormIncorrect(c, writeErr.WriteErrors)
				return
			}
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}
