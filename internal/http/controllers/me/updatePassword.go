package me

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			input       *forms.UpdateMePasswordForm
			dbConn      *mongo.Database
			user        *models.UserModel
			newPassword string
			err         error
		)

		if input, err = requests.GetUpdateMePasswordForm(c); err != nil {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if input.NewPassword != input.NewPasswordConfirm {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": "new password confirmation not match"})
			return
		}

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

		if hash.Check(newPassword, user.Password) {
			responses.Basic(c, http.StatusUnprocessableEntity, gin.H{"message": "your password is same as old one"})
		}

		if newPassword, err = hash.Make(input.NewPassword); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		user.Password = newPassword

		if err := repositories.UpdateUser(ctx, dbConn, user); err != nil {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		responses.Basic(c, http.StatusNoContent, nil)
	}
}
