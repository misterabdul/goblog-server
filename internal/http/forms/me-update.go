package forms

import (
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

type UpdateMeForm struct {
	FirstName string `json:"firstname" binding:"omitempty,max=50"`
	LastName  string `json:"lastname" binding:"omitempty,max=50"`
	Username  string `json:"username" binding:"omitempty,min=5,max=16"`
	Email     string `json:"email" binding:"email"`
}

func (form *UpdateMeForm) Validate(
	userService *service.UserService,
	me *models.UserModel,
) (err error) {
	if err = checkUpdateUsername(userService, form.Username, me); err != nil {
		return err
	}
	if err = checkUpdateEmail(userService, form.Email, me); err != nil {
		return err
	}

	return nil
}

func (form *UpdateMeForm) ToUserModel(
	me *models.UserModel,
) (updatedMe *models.UserModel) {
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
