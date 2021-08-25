package requests

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/forms"
)

func GetCreateCommentForm(c *gin.Context) (*forms.CreateCommentForm, error) {
	var createComment forms.CreateCommentForm
	err := shouldBind(c, &createComment)

	return &createComment, err
}

func GetReplyCommentForm(c *gin.Context) (*forms.ReplyCommmentForm, error) {
	var replyComment forms.ReplyCommmentForm
	err := shouldBind(c, &replyComment)

	return &replyComment, err
}
