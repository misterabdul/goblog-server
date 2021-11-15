package requests

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/forms"
)

func GetPasswordConfirmForm(c *gin.Context) (form *forms.PasswordConfirmForm, err error) {
	var _form forms.PasswordConfirmForm

	err = shouldBind(c, &_form)

	return &_form, err
}
