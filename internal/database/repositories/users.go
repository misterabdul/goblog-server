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

const userCollection = "users"

// Get single user
func ReadOneUser(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (user *models.UserModel, err error) {
	var (
		collection = dbConn.Collection(userCollection)
		_user      models.UserModel
	)

	if err = collection.FindOne(ctx, filter, opts...).Decode(&_user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_user, nil
}

// Get multiple users
func ReadManyUsers(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (users []*models.UserModel, err error) {
	var (
		collection = dbConn.Collection(userCollection)
		user       *models.UserModel
		cursor     *mongo.Cursor
	)

	if cursor, err = collection.Find(ctx, filter, opts...); err != nil {
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
func CountUsers(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {
	var collection = dbConn.Collection(userCollection)

	return collection.CountDocuments(
		ctx, filter, opts...)
}

// Save new user
func SaveOneUser(
	dbConn *mongo.Database,
	ctx context.Context,
	user *models.UserModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var (
		collection = dbConn.Collection(userCollection)
		insRes     *mongo.InsertOneResult
		insertedID primitive.ObjectID
		ok         bool
	)

	if insRes, err = collection.InsertOne(ctx, user, opts...); err != nil {
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
func UpdateOneUser(
	dbConn *mongo.Database,
	ctx context.Context,
	user *models.UserModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var collection = dbConn.Collection(userCollection)

	_, err = collection.UpdateOne(
		ctx, bson.M{"_id": user.UID}, bson.M{"$set": user}, opts...)

	return err
}

// Delete user
func DeleteOneUser(
	dbConn *mongo.Database,
	ctx context.Context,
	user *models.UserModel,
	opts ...*options.DeleteOptions,
) (err error) {
	var collection = dbConn.Collection(userCollection)

	_, err = collection.DeleteOne(
		ctx, bson.M{"_id": user.UID}, opts...)

	return err
}
