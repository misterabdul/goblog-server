package requests

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/forms"
)

func GetCreateCommentForm(c *gin.Context) (
	form *forms.CreateCommentForm,
	err error,
) {
	var _form = forms.CreateCommentForm{}

	err = shouldBind(c, &_form)

	return &_form, err
}

func GetCreateCommentReplyForm(c *gin.Context) (
	form *forms.CreateCommentReplyForm,
	err error,
) {
	var _form = forms.CreateCommentReplyForm{}

	err = shouldBind(c, &_form)

	return &_form, err
}
