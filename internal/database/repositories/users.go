package repositories

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/database/models"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(
	dbConn *mongo.Database,
) *UserRepository {

	return &UserRepository{
		collection: dbConn.Collection("users")}
}

// Get single user
func (r UserRepository) ReadOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (user *models.UserModel, err error) {
	var _user models.UserModel

	if err = r.collection.FindOne(
		ctx, filter, opts...,
	).Decode(&_user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_user, nil
}

// Get multiple users
func (r UserRepository) ReadMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (users []*models.UserModel, err error) {
	var (
		user   *models.UserModel
		cursor *mongo.Cursor
	)

	if cursor, err = r.collection.Find(
		ctx, filter, opts...,
	); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		user = &models.UserModel{}
		if err = cursor.Decode(user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// Count total users
func (r UserRepository) Count(
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {

	return r.collection.CountDocuments(
		ctx, filter, opts...,
	)
}

// Save new user
func (r UserRepository) Save(
	ctx context.Context,
	user *models.UserModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID primitive.ObjectID
		ok         bool
	)

	if insRes, err = r.collection.InsertOne(
		ctx, user,
	); err != nil {
		return err
	}
	if insertedID, ok = insRes.InsertedID.(primitive.ObjectID); !ok {
		return errors.New("unable to assert inserted uid")
	}
	if user.UID != insertedID {
		return errors.New("inserted uid is not same with database")
	}

	return nil
}

// Update user
func (r UserRepository) Update(
	ctx context.Context,
	user *models.UserModel,
) (err error) {
	_, err = r.collection.UpdateByID(
		ctx, user.UID, bson.M{"$set": user})

	return err
}

// Delete user
func (r UserRepository) Delete(
	ctx context.Context,
	user *models.UserModel,
) (err error) {
	_, err = r.collection.DeleteOne(
		ctx, bson.M{"_id": user.UID})

	return err
}
