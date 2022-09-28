package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/database/repositories"
)

type user struct {
	dbConn *mongo.Database
}

func newUserService(
	dbConn *mongo.Database,
) (service *user) {

	return &user{dbConn: dbConn}
}

// Get single user
func (s *user) GetOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (user *models.UserModel, err error) {

	return repositories.ReadOneUser(
		s.dbConn, ctx, filter, opts...)
}

// Get multiple users
func (s *user) GetMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (categories []*models.UserModel, err error) {

	return repositories.ReadManyUsers(
		s.dbConn, ctx, filter, opts...)
}

// Get total users count
func (s *user) Count(
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {

	return repositories.CountUsers(
		s.dbConn, ctx, filter, opts...)
}

// Create new user
func (s *user) SaveOne(
	ctx context.Context,
	user *models.UserModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	user.UID = primitive.NewObjectID()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.DeletedAt = nil

	return repositories.SaveOneUser(
		s.dbConn, ctx, user, opts...)
}

// Update user
func (s *user) UpdateOne(
	ctx context.Context,
	user *models.UserModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	user.UpdatedAt = now

	return repositories.UpdateOneUser(
		s.dbConn, ctx, user, opts...)
}

// Delete user to trash
func (s *user) TrashOne(
	ctx context.Context,
	user *models.UserModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	user.DeletedAt = now

	return repositories.UpdateOneUser(
		s.dbConn, ctx, user, opts...)
}

// Restore user from trash
func (s *user) RestoreOne(
	ctx context.Context,
	user *models.UserModel,
	opts ...*options.UpdateOptions,
) (err error) {
	user.DeletedAt = nil

	return repositories.UpdateOneUser(
		s.dbConn, ctx, user, opts...)
}

// Permanently delete user
func (s *user) DeleteOne(
	ctx context.Context,
	user *models.UserModel,
	opts ...*options.DeleteOptions,
) (err error) {

	return repositories.DeleteOneUser(
		s.dbConn, ctx, user, opts...)
}
