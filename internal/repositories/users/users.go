package users

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/models"
)

// Get the user collection
func getUserCollection(dbConn *mongo.Database) *mongo.Collection {
	return dbConn.Collection("users")
}

// Get single user
func GetUser(ctx context.Context, dbConn *mongo.Database, filter interface{}) (*models.UserModel, error) {
	var user models.UserModel
	if err := getUserCollection(dbConn).FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Get multiple users
func GetUsers(ctx context.Context, dbConn *mongo.Database, filter interface{}, show int, order string, asc bool) ([]*models.UserModel, error) {
	cursor, err := getUserCollection(dbConn).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.UserModel
	for cursor.Next(ctx) {
		var user models.UserModel
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// Create new user
func CreateUser(ctx context.Context, dbConn *mongo.Database, user *models.UserModel) error {
	now := primitive.NewDateTimeFromTime(time.Now())

	user.CreatedAt = now
	user.UpdatedAt = now
	user.DeletedAt = nil

	insRes, err := getUserCollection(dbConn).InsertOne(ctx, user)
	if err != nil {
		return err
	}
	insertedID, ok := insRes.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("unable to assert inserted uid")
	}
	user.UID = insertedID

	return nil
}

// Update user
func UpdateUser(ctx context.Context, dbConn *mongo.Database, user *models.UserModel) error {
	now := primitive.NewDateTimeFromTime(time.Now())

	user.UpdatedAt = now

	_, err := getUserCollection(dbConn).UpdateByID(ctx, user.UID, user)

	return err
}

// Mark user trash
func TrashUser(ctx context.Context, dbConn *mongo.Database, user *models.UserModel) error {
	now := primitive.NewDateTimeFromTime(time.Now())

	user.DeletedAt = now

	_, err := getUserCollection(dbConn).UpdateByID(ctx, user.UID, user)

	return err
}

// Permanently delete user
func DeleteUser(ctx context.Context, dbConn *mongo.Database, user *models.UserModel) error {
	_, err := getUserCollection(dbConn).DeleteOne(ctx, bson.M{"_id": user.UID})

	return err
}
