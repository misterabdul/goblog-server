package requests

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/forms"
)

func GetCreatePostForm(c *gin.Context) (form *forms.CreatePostForm, err error) {
	var _form forms.CreatePostForm

	err = shouldBind(c, &_form)

	return &_form, err
}

func GetUpdatePostForm(c *gin.Context) (form *forms.UpdatePostForm, err error) {
	var _form forms.UpdatePostForm

	err = shouldBind(c, &_form)

	return &_form, err
}
