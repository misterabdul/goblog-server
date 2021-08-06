package requests

import (
	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/forms"
)

func GetUpdateMeForm(c *gin.Context) (*forms.UpdateMeForm, error) {
	var updateMe forms.UpdateMeForm
	err := shouldBind(c, &updateMe)

	return &updateMe, err
}

func GetUpdateMePasswordForm(c *gin.Context) (*forms.UpdateMePasswordForm, error) {
	var updateMePassword forms.UpdateMePasswordForm
	err := shouldBind(c, &updateMePassword)

	return &updateMePassword, err
}
