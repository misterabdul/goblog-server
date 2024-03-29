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

type CreateUserForm struct {
	FirstName       string `json:"firstName" binding:"max=50"`
	LastName        string `json:"lastName" binding:"max=50"`
	Username        string `json:"username" binding:"required,alphanum,min=5,max=16"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8,max=32"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required,min=8,max=32"`
	Roles           []int  `json:"roles" binding:"omitempty,dive,number"`
}

func (form *CreateUserForm) Validate(
	svc *service.Service,
	ctx context.Context,
	creator *models.UserModel,
) (err error) {
	if err = isProperRoles(form.Roles); err != nil {
		return err
	}
	if strings.Compare(form.Password, form.PasswordConfirm) != 0 {
		return errors.New("password confirm not same")
	}
	if err = checkUsername(svc, ctx, form.Username); err != nil {
		return err
	}
	if err = checkEmail(svc, ctx, form.Email); err != nil {
		return err
	}

	return nil
}

func (form *CreateUserForm) ToUserModel() (user *models.UserModel, err error) {
	var (
		now      = primitive.NewDateTimeFromTime(time.Now())
		password string
	)

	if password, err = hash.Make(form.Password); err != nil {
		return nil, err
	}
	return &models.UserModel{
		UID:       primitive.NewObjectID(),
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Username:  form.Username,
		Email:     form.Email,
		Password:  password,
		Roles:     getRoles(form.Roles, now),
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}, nil
}

func isProperRoles(createdRoles []int) (err error) {
	var createdRoleLevel = 3

	for _, createdRole := range createdRoles {
		if createdRole < createdRoleLevel {
			createdRoleLevel = createdRole
		}
	}
	if createdRoleLevel <= 1 {
		return errors.New("cannot create admin role, you must adminize them later")
	}

	return nil
}

func checkUsername(
	svc *service.Service,
	ctx context.Context,
	formUsername string,
) (err error) {
	var users []*models.UserModel

	if users, err = svc.User.GetMany(ctx, bson.M{
		"username": bson.M{"$eq": formUsername},
	}); err != nil {
		return err
	}
	if len(users) > 0 {
		return errors.New("username exists")
	}

	return nil
}

func checkEmail(
	svc *service.Service,
	ctx context.Context,
	formEmail string,
) (err error) {
	var users []*models.UserModel

	if users, err = svc.User.GetMany(ctx, bson.M{
		"email": bson.M{"$eq": formEmail},
	}); err != nil {
		return err
	}
	if len(users) > 0 {
		return errors.New("email exists")
	}

	return nil
}

func getRoles(formRoles []int, now primitive.DateTime) (roles []models.UserRole) {
	var roleExists bool

	roles = []models.UserRole{}
	for _, roleNumber := range formRoles {
		roleExists = false
		for _, role := range roles {
			if roleNumber == role.Level {
				roleExists = true
				break
			}
		}
		if !roleExists {
			switch roleNumber {
			case 0:
				roles = append(roles, models.UserRole{
					Level: 0,
					Name:  "SuperAdmin",
					Since: now,
				})
			case 1:
				roles = append(roles, models.UserRole{
					Level: 1,
					Name:  "Admin",
					Since: now,
				})
			case 2:
				roles = append(roles, models.UserRole{
					Level: 2,
					Name:  "Editor",
					Since: now,
				})
			case 3:
				roles = append(roles, models.UserRole{
					Level: 3,
					Name:  "Writer",
					Since: now,
				})
			}
		}
	}

	return roles
}
