package requests

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/forms"
)

func GetCreateUserForm(c *gin.Context) (form *forms.CreateUserForm, err error) {
	var _form = forms.CreateUserForm{}

	err = shouldBind(c, &_form)

	return &_form, err
}

func GetUpdateUserForm(c *gin.Context) (form *forms.UpdateUserForm, err error) {
	var _form = forms.UpdateUserForm{}

	err = shouldBind(c, &_form)

	return &_form, err
}
