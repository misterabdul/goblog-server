package users

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/internal/requests"
	"github.com/misterabdul/goblog-server/internal/responses"
	"github.com/misterabdul/goblog-server/pkg/hash"
)

func SignIn(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input *models.SignInModel
	var dbConn *mongo.Database
	var user *models.UserModel
	var err error

	if input, err = requests.GetSignInModel(c); err != nil {
		responses.Basic(c, http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	if dbConn, err = repositories.GetDBConnDefault(ctx); err != nil {
		responses.Basic(c, http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if user, err = repositories.GetUser(ctx, dbConn, bson.M{"$or": []bson.M{
		{"username": input.Username},
		{"email": input.Username}}}); err != nil {
		responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "Wrong username or password."})
		return
	}
	if !hash.Check(input.Password, user.Password) {
		responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "Wrong username or password."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success."})
}

func SignInRefresh(c *gin.Context) {

	c.JSON(http.StatusNotImplemented, gin.H{})
}

func SignUp(c *gin.Context) {

	c.JSON(http.StatusNotImplemented, gin.H{})
}

func SignOut(c *gin.Context) {

	c.JSON(http.StatusNotImplemented, gin.H{})
}
