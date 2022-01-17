package repositories

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/models"
)

func getCommentCollection(
	dbConn *mongo.Database,
) (commentCollection *mongo.Collection) {
	return dbConn.Collection("comments")
}

// Get single comment
func GetComment(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.FindOneOptions,
) (comment *models.CommentModel, err error) {
	var _comment models.CommentModel
	if err = getCommentCollection(dbConn).FindOne(
		ctx, filter, opts...,
	).
		Decode(&_comment); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_comment, nil
}

// Get multiple comments
func GetComments(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.FindOptions,
) (comments []*models.CommentModel, err error) {
	var (
		comment *models.CommentModel
		cursor  *mongo.Cursor
	)

	if cursor, err = getCommentCollection(dbConn).Find(
		ctx, filter, opts...,
	); err != nil {
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

// Save new comment
func SaveComment(
	ctx context.Context,
	dbConn *mongo.Database,
	comment *models.CommentModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID primitive.ObjectID
		ok         bool
	)

	if insRes, err = getCommentCollection(dbConn).InsertOne(
		ctx, comment,
	); err != nil {
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
func UpdateComment(
	ctx context.Context,
	dbConn *mongo.Database,
	comment *models.CommentModel,
) (err error) {
	_, err = getCommentCollection(dbConn).UpdateByID(
		ctx, comment.UID, bson.M{"$set": comment})

	return err
}

// Delete comment
func DeleteComment(
	ctx context.Context,
	dbConn *mongo.Database,
	comment *models.CommentModel,
) (err error) {
	_, err = getCommentCollection(dbConn).DeleteOne(
		ctx, bson.M{"_id": comment.UID})

	return err
}
