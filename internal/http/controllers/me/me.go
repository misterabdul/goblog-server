package me

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
)

func GetMe(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me  *models.UserModel
			err error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}

		responses.Me(c, me)
	}
}
