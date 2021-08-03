package users

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

func GetMe(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var (
		dbConn *mongo.Database
		user   *models.UserModel
		err    error
	)

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

	responses.Me(c, user)
}

func UpdateMe(c *gin.Context) {

	c.JSON(http.StatusNotImplemented, gin.H{})
}

func UpdateMePassword(c *gin.Context) {

	c.JSON(http.StatusNotImplemented, gin.H{})
}
