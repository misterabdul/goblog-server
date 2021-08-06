package requests

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/forms"
)

func GetPasswordConfirmForm(c *gin.Context) (*forms.PasswordConfirmForm, error) {
	var passwordConfirmForm forms.PasswordConfirmForm
	err := shouldBind(c, &passwordConfirmForm)

	return &passwordConfirmForm, err
}
