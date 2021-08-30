package repositories

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/models"
)

func getCommentCollection(dbConn *mongo.Database) *mongo.Collection {
	return dbConn.Collection("comments")
}

// Get single comment
func GetComment(ctx context.Context, dbConn *mongo.Database, filter interface{}) (*models.CommentModel, error) {
	var comment models.CommentModel
	if err := getCommentCollection(dbConn).FindOne(ctx, filter).Decode(&comment); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &comment, nil
}

// Get multiple comments
func GetComments(ctx context.Context, dbConn *mongo.Database, filter interface{}, show int, order string, asc bool) ([]*models.CommentModel, error) {
	cursor, err := getCommentCollection(dbConn).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []*models.CommentModel
	for cursor.Next(ctx) {
		var comment models.CommentModel
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}

// Create new comment
func CreateComment(ctx context.Context, dbConn *mongo.Database, comment *models.CommentModel) error {
	now := primitive.NewDateTimeFromTime(time.Now())

	comment.UID = primitive.NewObjectID()
	comment.CreatedAt = now
	comment.DeletedAt = nil

	insRes, err := getCommentCollection(dbConn).InsertOne(ctx, comment)
	if err != nil {
		return err
	}
	insertedID, ok := insRes.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("unable to assert inserted uid")
	}
	if comment.UID != insertedID {
		return errors.New("inserted uid is not same with database")
	}

	return nil
}

// Mark comment trash
func TrashComment(ctx context.Context, dbConn *mongo.Database, comment *models.CommentModel) error {
	now := primitive.NewDateTimeFromTime(time.Now())

	comment.DeletedAt = now

	_, err := getCommentCollection(dbConn).UpdateByID(ctx, comment.UID, bson.M{"$set": comment})

	return err
}

// Unmark the trash from comment
func DetrashComment(ctx context.Context, dbConn *mongo.Database, comment *models.CommentModel) error {
	comment.DeletedAt = nil

	_, err := getCommentCollection(dbConn).UpdateByID(ctx, comment.UID, bson.M{"$set": comment})

	return err
}

// Permanently delete comment
func DeleteComment(ctx context.Context, dbConn *mongo.Database, comment *models.CommentModel) error {
	_, err := getCommentCollection(dbConn).DeleteOne(ctx, bson.M{"_id": comment.UID})

	return err
}
