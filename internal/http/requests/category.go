package requests

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/forms"
)

func GetCreateCategoryForm(c *gin.Context) (form *forms.CreateCategoryForm, err error) {
	var _form = forms.CreateCategoryForm{}

	err = shouldBind(c, &_form)

	return &_form, err
}

func GetUpdateCategoryForm(c *gin.Context) (form *forms.UpdateCategoryForm, err error) {
	var _form = forms.UpdateCategoryForm{}

	err = shouldBind(c, &_form)

	return &_form, err
}
