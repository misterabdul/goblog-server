package forms

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/service"
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
	svc *service.Service,
	ctx context.Context,
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
		if err = checkUpdateUsername(svc, ctx, form.Username, target); err != nil {
			return err
		}
	}
	if strings.Compare(form.Email, target.Email) != 0 {
		if err = checkUpdateEmail(svc, ctx, form.Email, target); err != nil {
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

func checkUpdateUsername(
	svc *service.Service,
	ctx context.Context,
	formUsername string,
	target *models.UserModel,
) (err error) {
	var users []*models.UserModel

	if users, err = svc.User.GetMany(ctx, bson.M{
		"$and": []bson.M{
			{"_id": bson.M{"$ne": target.UID}},
			{"username": bson.M{"$eq": formUsername}}},
	}); err != nil {
		return err
	}
	if len(users) > 0 {
		return errors.New("username exists")
	}

	return nil
}

func checkUpdateEmail(
	svc *service.Service,
	ctx context.Context,
	formEmail string,
	target *models.UserModel,
) (err error) {
	var users []*models.UserModel

	if users, err = svc.User.GetMany(ctx, bson.M{
		"$and": []bson.M{
			{"_id": bson.M{"$ne": target.UID}},
			{"email": bson.M{"$eq": formEmail}}},
	}); err != nil {
		return err
	}
	if len(users) > 0 {
		return errors.New("email exists")
	}

	return nil
}
