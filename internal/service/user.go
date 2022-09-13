package service

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/database/repositories"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
)

type UserService struct {
	c          *gin.Context
	ctx        context.Context
	dbConn     *mongo.Database
	repository *repositories.UserRepository
}

func NewUserService(
	c *gin.Context,
	ctx context.Context,
	dbConn *mongo.Database,
) *UserService {

	return &UserService{
		c:          c,
		ctx:        ctx,
		dbConn:     dbConn,
		repository: repositories.NewUserRepository(dbConn)}
}

// Get single user
func (s *UserService) GetUser(filter interface{}) (
	user *models.UserModel,
	err error,
) {

	return s.repository.ReadOne(
		s.ctx, filter)
}

// Get multiple users
func (s *UserService) GetUsers(filter interface{}) (
	categories []*models.UserModel,
	err error,
) {

	return s.repository.ReadMany(
		s.ctx, filter,
		internalGin.GetFindOptions(s.c))
}

// Get total users count
func (s *UserService) GetUserCount(filter interface{}) (
	count int64, err error,
) {

	return s.repository.Count(
		s.ctx, filter,
		internalGin.GetCountOptions(s.c))
}

// Create new user
func (s *UserService) CreateUser(
	user *models.UserModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	user.UID = primitive.NewObjectID()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.DeletedAt = nil

	return s.repository.Save(
		s.ctx, user)
}

// Update user
func (s *UserService) UpdateUser(
	user *models.UserModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	user.UpdatedAt = now

	return s.repository.Update(
		s.ctx, user)
}

// Delete user to trash
func (s *UserService) TrashUser(
	user *models.UserModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	user.DeletedAt = now

	return s.repository.Update(
		s.ctx, user)
}

// Restore user from trash
func (s *UserService) DetrashUser(
	user *models.UserModel,
) (err error) {
	user.DeletedAt = nil

	return s.repository.Update(
		s.ctx, user)
}

// Permanently delete user
func (s *UserService) DeleteUser(
	user *models.UserModel,
) (err error) {

	return s.repository.Delete(
		s.ctx, user)
}
