package requests

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/forms"
)

func GetUpdateMeForm(c *gin.Context) (form *forms.UpdateMeForm, err error) {
	var (
		_form = forms.UpdateMeForm{}
	)

	err = shouldBind(c, &_form)

	return &_form, err
}

func GetUpdateMePasswordForm(c *gin.Context) (form *forms.UpdateMePasswordForm, err error) {
	var (
		_form = forms.UpdateMePasswordForm{}
	)

	err = shouldBind(c, &_form)

	return &_form, err
}
