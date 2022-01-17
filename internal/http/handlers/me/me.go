package me

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
)

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
