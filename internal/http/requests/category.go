package requests

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/forms"
)

func GetCreateCategoryForm(c *gin.Context) (*forms.CreateCategoryForm, error) {
	var createCategory forms.CreateCategoryForm
	err := shouldBind(c, &createCategory)

	return &createCategory, err
}

func GetUpdateCategoryForm(c *gin.Context) (*forms.UpdateCategoryForm, error) {
	var updateCategory forms.UpdateCategoryForm
	err := shouldBind(c, &updateCategory)

	return &updateCategory, err
}
