package repositories

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/models"
)

// Get the user collection
func getUserCollection(dbConn *mongo.Database) *mongo.Collection {
	return dbConn.Collection("users")
}

// Get single user
func GetUser(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.FindOneOptions,
) (user *models.UserModel, err error) {
	var _user models.UserModel
	if err = getUserCollection(dbConn).FindOne(ctx, filter, opts...).
		Decode(&_user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_user, nil
}

// Get multiple users
func GetUsers(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.FindOptions,
) (users []*models.UserModel, err error) {
	var (
		user   *models.UserModel
		cursor *mongo.Cursor
	)

	if cursor, err = getUserCollection(dbConn).Find(ctx, filter, opts...); err != nil {
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

// Create new user
func CreateUser(
	ctx context.Context,
	dbConn *mongo.Database,
	user *models.UserModel,
) (err error) {
	var (
		now        = primitive.NewDateTimeFromTime(time.Now())
		insRes     *mongo.InsertOneResult
		insertedID primitive.ObjectID
		ok         bool
	)

	user.UID = primitive.NewObjectID()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.DeletedAt = nil
	if insRes, err = getUserCollection(dbConn).InsertOne(ctx, user); err != nil {
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
func UpdateUser(
	ctx context.Context,
	dbConn *mongo.Database,
	user *models.UserModel,
) (err error) {
	now := primitive.NewDateTimeFromTime(time.Now())
	user.UpdatedAt = now
	_, err = getUserCollection(dbConn).
		UpdateByID(ctx, user.UID, bson.M{"$set": user})

	return err
}

// Mark user trash
func TrashUser(
	ctx context.Context,
	dbConn *mongo.Database,
	user *models.UserModel,
) (err error) {
	now := primitive.NewDateTimeFromTime(time.Now())
	user.DeletedAt = now
	_, err = getUserCollection(dbConn).
		UpdateByID(ctx, user.UID, bson.M{"$set": user})

	return err
}

func DetrashUser(
	ctx context.Context,
	dbConn *mongo.Database,
	user *models.UserModel,
) (err error) {
	user.DeletedAt = nil
	_, err = getUserCollection(dbConn).
		UpdateByID(ctx, user.UID, bson.M{"$set": user})

	return err
}

// Permanently delete user
func DeleteUser(
	ctx context.Context,
	dbConn *mongo.Database,
	user *models.UserModel,
) (err error) {
	_, err = getUserCollection(dbConn).
		DeleteOne(ctx, bson.M{"_id": user.UID})

	return err
}
