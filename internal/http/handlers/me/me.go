package me

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
)

// @Tags        Me
// @Summary     Get Me
// @Description Get my user data.
// @Router      /v1/auth/me [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Success     200 {object} object{data=object{uid=string,username=string,email=string,firstName=string,lastName=string,roles=[]object{level=int,name=string,since=time},updatedAt=time,createdAt=time}}
// @Failure     401 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func GetMe(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			me  *models.UserModel
			err error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}

		responses.Me(c, me)
	}
}
