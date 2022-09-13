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

type CommentRepository struct {
	collection *mongo.Collection
}

func NewCommentRepository(
	dbConn *mongo.Database,
) *CommentRepository {

	return &CommentRepository{
		collection: dbConn.Collection("comments")}
}

// Get single comment
func (r CommentRepository) ReadOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (comment *models.CommentModel, err error) {
	var _comment models.CommentModel

	if err = r.collection.FindOne(
		ctx, filter, opts...,
	).Decode(&_comment); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_comment, nil
}

// Get multiple comments
func (r CommentRepository) ReadMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (comments []*models.CommentModel, err error) {
	var (
		comment *models.CommentModel
		cursor  *mongo.Cursor
	)

	if cursor, err = r.collection.Find(
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

// Count total posts
func (r CommentRepository) Count(
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {

	return r.collection.CountDocuments(
		ctx, filter, opts...,
	)
}

// Save new comment
func (r CommentRepository) Save(
	ctx context.Context,
	comment *models.CommentModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID primitive.ObjectID
		ok         bool
	)

	if insRes, err = r.collection.InsertOne(
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
func (r CommentRepository) Update(
	ctx context.Context,
	comment *models.CommentModel,
) (err error) {
	_, err = r.collection.UpdateByID(
		ctx, comment.UID, bson.M{"$set": comment})

	return err
}

// Delete comment
func (r CommentRepository) Delete(
	ctx context.Context,
	comment *models.CommentModel,
) (err error) {
	_, err = r.collection.DeleteOne(
		ctx, bson.M{"_id": comment.UID})

	return err
}
