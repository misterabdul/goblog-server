package requests

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/forms"
)

func GetCreatePageForm(c *gin.Context) (form *forms.CreatePageForm, err error) {
	var _form = forms.CreatePageForm{}

	err = shouldBind(c, &_form)

	return &_form, err
}

func GetUpdatePageForm(c *gin.Context) (form *forms.UpdatePageForm, err error) {
	var _form = forms.UpdatePageForm{}

	err = shouldBind(c, &_form)

	return &_form, err
}
