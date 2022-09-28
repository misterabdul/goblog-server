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

const commentCollection = "comments"

// Get single comment
func ReadOneComment(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (comment *models.CommentModel, err error) {
	var (
		collection = dbConn.Collection(commentCollection)
		_comment   models.CommentModel
	)

	if err = collection.FindOne(ctx, filter, opts...).Decode(&_comment); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_comment, nil
}

// Get multiple comments
func ReadManyComments(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (comments []*models.CommentModel, err error) {
	var (
		collection = dbConn.Collection(commentCollection)
		comment    *models.CommentModel
		cursor     *mongo.Cursor
	)

	if cursor, err = collection.Find(ctx, filter, opts...); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		comment = &models.CommentModel{}
		if err = cursor.Decode(comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// Count total posts
func CountComments(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {
	var collection = dbConn.Collection(commentCollection)

	return collection.CountDocuments(
		ctx, filter, opts...)
}

// Save new comment
func SaveOneComment(
	dbConn *mongo.Database,
	ctx context.Context,
	comment *models.CommentModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var (
		collection = dbConn.Collection(commentCollection)
		insRes     *mongo.InsertOneResult
		insertedID primitive.ObjectID
		ok         bool
	)

	if insRes, err = collection.InsertOne(ctx, comment, opts...); err != nil {
		return err
	}
	if insertedID, ok = insRes.InsertedID.(primitive.ObjectID); !ok {
		return errors.New("unable to assert inserted uid")
	}
	if comment.UID != insertedID {
		return errors.New("inserted uid is not same with database")
	}

	return nil
}

// Update comment
func UpdateOneComment(
	dbConn *mongo.Database,
	ctx context.Context,
	comment *models.CommentModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var collection = dbConn.Collection(commentCollection)

	_, err = collection.UpdateOne(
		ctx, bson.M{"_id": comment.UID}, bson.M{"$set": comment}, opts...)

	return err
}

// Delete comment
func DeleteOneComment(
	dbConn *mongo.Database,
	ctx context.Context,
	comment *models.CommentModel,
	opts ...*options.DeleteOptions,
) (err error) {
	var collection = dbConn.Collection(commentCollection)

	_, err = collection.DeleteOne(
		ctx, bson.M{"_id": comment.UID}, opts...)

	return err
}
