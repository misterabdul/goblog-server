package forms

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/pkg/hash"
)

type CreateUserForm struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
	Roles           []int  `json:"roles"`
}

type UpdateUserForm struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
	Roles           []int  `json:"roles"`
}

func (form *CreateUserForm) Validate(
	ctx context.Context,
	dbConn *mongo.Database,
	creator *models.UserModel,
) (err error) {
	var (
		users []*models.UserModel
	)

	if err = isProperRoles(form.Roles); err != nil {
		return err
	}
	if strings.Compare(form.Password, form.PasswordConfirm) != 0 {
		return errors.New("password confirm not same")
	}
	if users, err = repositories.GetUsers(ctx, dbConn, bson.M{
		"username": bson.M{"$eq": form.Username},
	}); err != nil {
		return err
	}
	if len(users) > 0 {
		return errors.New("username exists")
	}
	if users, err = repositories.GetUsers(ctx, dbConn, bson.M{
		"email": bson.M{"$eq": form.Email},
	}); err != nil {
		return err
	}
	if len(users) > 0 {
		return errors.New("email exists")
	}

	return nil
}

func (form *CreateUserForm) ToUserModel() (user *models.UserModel, err error) {
	var (
		now      = primitive.NewDateTimeFromTime(time.Now())
		userId   = primitive.NewObjectID()
		password string
	)

	if password, err = hash.Make(form.Password); err != nil {
		return nil, err
	}
	return &models.UserModel{
		UID:       userId,
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

func (form *UpdateUserForm) Validate(
	ctx context.Context,
	dbConn *mongo.Database,
	creator *models.UserModel,
	target *models.UserModel,
) (err error) {
	var (
		users []*models.UserModel
	)

	if err = isProperRoles(form.Roles); err != nil {
		return err
	}
	if strings.Compare(form.Password, form.PasswordConfirm) != 0 {
		return errors.New("password confirm not same")
	}
	if strings.Compare(form.Username, target.Username) != 0 {
		if users, err = repositories.GetUsers(ctx, dbConn, bson.M{
			"username": bson.M{"$eq": form.Username},
		}); err != nil {
			return err
		}
		if len(users) > 0 {
			return errors.New("username exists")
		}
	}
	if strings.Compare(form.Email, target.Email) != 0 {
		if users, err = repositories.GetUsers(ctx, dbConn, bson.M{
			"email": bson.M{"$eq": form.Email},
		}); err != nil {
			return err
		}
		if len(users) > 0 {
			return errors.New("email exists")
		}
	}

	return nil
}

func (form *UpdateUserForm) ToUserModel(
	currentUser *models.UserModel,
) (updatedUser *models.UserModel, err error) {
	var (
		now      = primitive.NewDateTimeFromTime(time.Now())
		password string
	)

	if len(form.FirstName) > 0 {
		currentUser.FirstName = form.FirstName
	}
	if len(form.LastName) > 0 {
		currentUser.LastName = form.LastName
	}
	if len(form.Username) > 0 {
		currentUser.Username = form.Username
	}
	if len(form.Email) > 0 {
		currentUser.Email = form.Email
	}
	if len(form.Password) > 0 {
		if password, err = hash.Make(form.Password); err != nil {
			return nil, err
		}
		currentUser.Password = password
	}
	if len(form.Roles) > 0 {
		currentUser.Roles = getRoles(form.Roles, now)
	}

	return currentUser, nil
}

func isProperRoles(createdRoles []int) (err error) {
	var (
		createdRoleLevel = 3
	)

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

func getRoles(formRoles []int, now primitive.DateTime) (roles []models.UserRoles) {
	var (
		roleExists bool
	)

	roles = []models.UserRoles{}
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
				roles = append(roles, models.UserRoles{
					Level: 0,
					Name:  "SuperAdmin",
					Since: now,
				})
			case 1:
				roles = append(roles, models.UserRoles{
					Level: 1,
					Name:  "Admin",
					Since: now,
				})
			case 2:
				roles = append(roles, models.UserRoles{
					Level: 2,
					Name:  "Editor",
					Since: now,
				})
			case 3:
				roles = append(roles, models.UserRoles{
					Level: 3,
					Name:  "Writer",
					Since: now,
				})
			}
		}
	}

	return roles
}
