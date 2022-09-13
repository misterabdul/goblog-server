package forms

import (
	"errors"
	"strings"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/pkg/hash"
)

type UpdateMePasswordForm struct {
	NewPassword        string `json:"newpassword" bind:"required,min=8,max=32"`
	NewPasswordConfirm string `json:"newpasswordconfirm" bind:"required,min=8,max=32"`
}

func (form *UpdateMePasswordForm) Validate(me *models.UserModel) (err error) {
	if strings.Compare(form.NewPassword, form.NewPasswordConfirm) != 0 {
		return errors.New("new password confirm not same")
	}
	if isSamePassowrd := hash.Check(form.NewPassword, me.Password); isSamePassowrd {
		return errors.New("new password is same as old password")
	}

	return
}

func (form *UpdateMePasswordForm) ToUserModel(
	me *models.UserModel,
) (updatedMe *models.UserModel, err error) {
	var newPassword string

	if newPassword, err = hash.Make(form.NewPassword); err != nil {
		return nil, err
	}
	me.Password = newPassword

	return me, nil
}
