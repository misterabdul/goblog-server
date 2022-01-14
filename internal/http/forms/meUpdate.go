package forms

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

type UpdateMeForm struct {
	FirstName string `json:"firstname" binding:"omitempty,max=50"`
	LastName  string `json:"lastname" binding:"omitempty,max=50"`
	Username  string `json:"username" binding:"omitempty,min=5,max=16"`
	Email     string `json:"email" binding:"email"`
}

func (form *UpdateMeForm) Validate(
	ctx context.Context,
	dbConn *mongo.Database,
	me *models.UserModel,
) (err error) {
	if err = checkUsername(ctx, dbConn, form.Username); err != nil {
		return err
	}
	if err = checkEmail(ctx, dbConn, form.Email); err != nil {
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

func checkUsername(
	ctx context.Context,
	dbConn *mongo.Database,
	formUsername string,
) (err error) {
	var (
		users []*models.UserModel
	)

	if users, err = repositories.GetUsers(ctx, dbConn, bson.M{
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
	ctx context.Context,
	dbConn *mongo.Database,
	formEmail string,
) (err error) {
	var (
		users []*models.UserModel
	)

	if users, err = repositories.GetUsers(ctx, dbConn, bson.M{
		"email": bson.M{"$eq": formEmail},
	}); err != nil {
		return err
	}
	if len(users) > 0 {
		return errors.New("email exists")
	}

	return nil
}
