package me

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/pkg/hash"
)

func UpdateMePassword(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn      *mongo.Database
			me          *models.UserModel
			form        *forms.UpdateMePasswordForm
			newPassword string
			err         error
		)

		if form, err = requests.GetUpdateMePasswordForm(c); err != nil {
			responses.Basic(c, http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}
		if form.NewPassword != form.NewPasswordConfirm {
			responses.Basic(c, http.StatusUnprocessableEntity, gin.H{"message": "new password confirmation not match"})
			return
		}
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if hash.Check(newPassword, me.Password) {
			responses.Basic(c, http.StatusUnprocessableEntity, gin.H{"message": "your password is same as old one"})
		}
		if newPassword, err = hash.Make(form.NewPassword); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		me.Password = newPassword
		if err := repositories.UpdateUser(ctx, dbConn, me); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responses.Basic(c, http.StatusNoContent, nil)
	}
}
