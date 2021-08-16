package requests

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/forms"
)

func GetCreatePostForm(c *gin.Context) (*forms.CreatePostForm, error) {
	var createPost forms.CreatePostForm
	err := shouldBind(c, &createPost)

	return &createPost, err
}

func GetUpdatePostForm(c *gin.Context) (*forms.UpdatePostForm, error) {
	var updatePost forms.UpdatePostForm
	err := shouldBind(c, &updatePost)

	return &updatePost, err
}
