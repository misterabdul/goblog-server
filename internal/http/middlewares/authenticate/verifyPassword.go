package authenticate

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
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/pkg/hash"
)

func VerifyPassword(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			input  *forms.PasswordConfirmForm
			dbConn *mongo.Database
			user   *models.UserModel
			err    error
		)

		if input, err = requests.GetPasswordConfirmForm(c); err != nil || input.Password == "" {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": "you must provide your password"})
			c.Abort()
			return
		}

		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		userUid, err := primitive.ObjectIDFromHex(c.GetString(AuthenticatedUserUid))
		if err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			c.Abort()
			return
		}

		if user, err = repositories.GetUser(ctx, dbConn, bson.M{"_id": userUid}); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			c.Abort()
			return
		}

		if !hash.Check(input.Password, user.Password) {
			responses.Basic(c, http.StatusBadRequest, gin.H{"message": "your password is wrong"})
			c.Abort()
			return
		}

		c.Next()
	}
}
