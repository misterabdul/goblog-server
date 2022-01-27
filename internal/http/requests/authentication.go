package requests

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/forms"
)

func GetSignInForm(c *gin.Context) (form *forms.SignInForm, err error) {
	var _form = forms.SignInForm{}

	err = shouldBind(c, &_form)

	return &_form, err
}
