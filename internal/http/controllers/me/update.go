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
)

func UpdateMe(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn *mongo.Database
			me     *models.UserModel
			form   *forms.UpdateMeForm
			err    error
		)

		if form, err = requests.GetUpdateMeForm(c); err != nil {
			responses.Basic(c, http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
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

		if err := repositories.UpdateUser(ctx, dbConn, forms.UpdateMeUserModel(form, me)); err != nil {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		responses.Basic(c, http.StatusNoContent, nil)
	}
}
