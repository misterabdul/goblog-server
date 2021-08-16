package forms

import "github.com/misterabdul/goblog-server/internal/models"

type UpdateMeForm struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`
	Email     string `json:"email" binding:"email"`
}

type UpdateMePasswordForm struct {
	NewPassword        string `json:"newpassword" bind:"required"`
	NewPasswordConfirm string `json:"newpasswordconfirm" bind:"required"`
}

func UpdateMeUserModel(form *UpdateMeForm, me *models.UserModel) *models.UserModel {
	if len(form.FirstName) > 0 {
		me.FirstName = form.FirstName
	}
	if len(form.LastName) > 0 {
		me.LastName = form.LastName
	}
	if len(form.Username) > 0 {
		me.Username = form.Username
	}
	if len(form.Email) > 0 {
		me.Email = form.Email
	}

	return me
}
