package requests

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/forms"
)

func GetSignInForm(c *gin.Context) (*forms.SignInForm, error) {
	var signIn forms.SignInForm
	err := shouldBind(c, &signIn)

	return &signIn, err
}
