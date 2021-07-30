package requests

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/models"
)

func GetSignInModel(c *gin.Context) (*models.SignInModel, error) {
	var signIn models.SignInModel
	if err := c.ShouldBind(&signIn); err != nil {
		return nil, err
	}
	return &signIn, nil
}
