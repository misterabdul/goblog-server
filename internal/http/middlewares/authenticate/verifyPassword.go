package authenticate

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/service"
	"github.com/misterabdul/goblog-server/pkg/hash"
)

func VerifyPassword(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			_, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			input     *forms.PasswordConfirmForm
			user      *models.UserModel
			err       error
		)

		defer cancel()
		if input, err = requests.GetPasswordConfirmForm(c); err != nil || input.Password == "" {
			responses.Basic(c, http.StatusUnprocessableEntity, gin.H{
				"message": "you must provide your password"})
			c.Abort()
			return
		}
		if user, err = GetAuthenticatedUser(c); err != nil {
			responses.Basic(c, http.StatusUnauthorized, gin.H{
				"message": "user not found"})
			c.Abort()
			return
		}
		if !hash.Check(input.Password, user.Password) {
			responses.Basic(c, http.StatusBadRequest, gin.H{
				"message": "your password is wrong"})
			c.Abort()
			return
		}
		c.Next()
	}
}
