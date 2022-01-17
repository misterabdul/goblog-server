package service

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/models"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

// Get single user
func (service *Service) GetUser(filter interface{}) (
	user *models.UserModel,
	err error,
) {

	return repositories.GetUser(
		service.ctx,
		service.dbConn,
		filter)
}

// Get multiple users
func (service *Service) GetUsers(filter interface{}) (
	categories []*models.UserModel,
	err error,
) {

	return repositories.GetUsers(
		service.ctx,
		service.dbConn,
		filter,
		internalGin.GetFindOptions(service.c))
}

// Create new user
func (service *Service) CreateUser(
	user *models.UserModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	user.UID = primitive.NewObjectID()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.DeletedAt = nil

	return repositories.SaveUser(
		service.ctx,
		service.dbConn,
		user)
}

// Update user
func (service *Service) UpdateUser(
	user *models.UserModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	user.UpdatedAt = now

	return repositories.UpdateUser(
		service.ctx,
		service.dbConn,
		user)
}

// Delete user to trash
func (service *Service) TrashUser(
	user *models.UserModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	user.DeletedAt = now

	return repositories.UpdateUser(
		service.ctx,
		service.dbConn,
		user)
}

// Restore user from trash
func (service *Service) DetrashUser(
	user *models.UserModel,
) (err error) {
	user.DeletedAt = nil

	return repositories.UpdateUser(
		service.ctx,
		service.dbConn,
		user)
}

// Permanently delete user
func (service *Service) DeleteUser(
	user *models.UserModel,
) (err error) {

	return repositories.DeleteUser(
		service.ctx,
		service.dbConn,
		user)
}
