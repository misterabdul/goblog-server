package forms

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/pkg/hash"
)

type UpdateUserForm struct {
	FirstName       string `json:"firstName" binding:"omitempty,max=50"`
	LastName        string `json:"lastName" binding:"omitempty,max=50"`
	Username        string `json:"username" binding:"omitempty,alphanum,min=5,max=16"`
	Email           string `json:"email" binding:"omitempty,email"`
	Password        string `json:"password" binding:"omitempty,min=8,max=32"`
	PasswordConfirm string `json:"passwordConfirm" binding:"omitempty,min=8,max=32"`
	Roles           []int  `json:"roles" binding:"omitempty,dive,number"`
}

func (form *UpdateUserForm) Validate(
	ctx context.Context,
	dbConn *mongo.Database,
	creator *models.UserModel,
	target *models.UserModel,
) (err error) {
	if err = isProperRoles(form.Roles); err != nil {
		return err
	}
	if strings.Compare(form.Password, form.PasswordConfirm) != 0 {
		return errors.New("password confirm not same")
	}
	if strings.Compare(form.Username, target.Username) != 0 {
		if err = checkUsername(ctx, dbConn, form.Username); err != nil {
			return err
		}
	}
	if strings.Compare(form.Email, target.Email) != 0 {
		if err = checkEmail(ctx, dbConn, form.Email); err != nil {
			return err
		}
	}

	return nil
}

func (form *UpdateUserForm) ToUserModel(
	user *models.UserModel,
) (updatedUser *models.UserModel, err error) {
	var (
		now      = primitive.NewDateTimeFromTime(time.Now())
		password string
	)

	if len(form.FirstName) > 0 {
		user.FirstName = form.FirstName
	}
	if len(form.LastName) > 0 {
		user.LastName = form.LastName
	}
	if len(form.Username) > 0 {
		user.Username = form.Username
	}
	if len(form.Email) > 0 {
		user.Email = form.Email
	}
	if len(form.Password) > 0 {
		if password, err = hash.Make(form.Password); err != nil {
			return nil, err
		}
		user.Password = password
	}
	if len(form.Roles) > 0 {
		user.Roles = getRoles(form.Roles, now)
	}

	return user, nil
}
