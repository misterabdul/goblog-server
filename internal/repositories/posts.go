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

type PostRepository struct {
	collection *mongo.Collection
}

type PostContentRepository struct {
	collection *mongo.Collection
}

func NewPostRepository(
	dbConn *mongo.Database,
) *PostRepository {

	return &PostRepository{
		collection: dbConn.Collection("posts")}
}

func NewPostContentRepository(
	dbConn *mongo.Database,
) *PostContentRepository {

	return &PostContentRepository{
		collection: dbConn.Collection("postContents")}
}

// Get single post
func (r PostRepository) ReadOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (post *models.PostModel, err error) {
	var _post models.PostModel

	if err = r.collection.FindOne(
		ctx, filter, opts...,
	).Decode(&_post); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_post, nil
}

// Get single post content
func (r PostContentRepository) ReadOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (
	postContent *models.PostContentModel,
	err error,
) {
	var _postContent models.PostContentModel

	if err = r.collection.FindOne(
		ctx, filter, opts...,
	).Decode(&_postContent); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_postContent, nil
}

// Get multiple posts
func (r PostRepository) ReadMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (posts []*models.PostModel, err error) {
	var (
		post   *models.PostModel
		cursor *mongo.Cursor
	)

	if cursor, err = r.collection.Find(
		ctx, filter, opts...,
	); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		post = &models.PostModel{}
		if err = cursor.Decode(post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// Count total posts
func (r PostRepository) Count(
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {

	return r.collection.CountDocuments(
		ctx, filter, opts...,
	)
}

// Save new post
func (r PostRepository) Save(
	ctx context.Context,
	post *models.PostModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID interface{}
		ok         bool
	)

	if insRes, err = r.collection.InsertOne(
		ctx, post,
	); err != nil {
		return err
	}
	if insertedID, ok = insRes.InsertedID.(primitive.ObjectID); !ok {
		return errors.New("unable to assert inserted uid")
	}
	if post.UID != insertedID {
		return errors.New("inserted uid is not same with database")
	}

	return nil
}

// Save new post content
func (r PostContentRepository) Save(
	ctx context.Context,
	postContent *models.PostContentModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID interface{}
		ok         bool
	)

	if insRes, err = r.collection.InsertOne(
		ctx, postContent,
	); err != nil {
		return err
	}
	if insertedID, ok = insRes.InsertedID.(primitive.ObjectID); !ok {
		return errors.New("unable to assert inserted uid")
	}
	if postContent.UID != insertedID {
		return errors.New("inserted uid is not same with database")
	}

	return nil
}

// Update post
func (r PostRepository) Update(
	ctx context.Context,
	post *models.PostModel,
) (err error) {
	if _, err = r.collection.UpdateByID(
		ctx, post.UID, bson.M{"$set": post},
	); err != nil {
		return err
	}

	return nil
}

// Update post content
func (r PostContentRepository) Update(
	ctx context.Context,
	postContent *models.PostContentModel,
) (err error) {
	if _, err = r.collection.UpdateByID(
		ctx, postContent.UID, bson.M{"$set": postContent},
	); err != nil {
		return err
	}

	return nil
}

// Delete post
func (r PostRepository) Delete(
	ctx context.Context,
	post *models.PostModel,
) (err error) {
	if _, err = r.collection.DeleteOne(
		ctx, bson.M{"_id": post.UID},
	); err != nil {
		return err
	}

	return nil
}

// Delete post content
func (r PostContentRepository) Delete(
	ctx context.Context,
	postContent *models.PostContentModel,
) (err error) {
	if _, err = r.collection.DeleteOne(
		ctx, bson.M{"_id": postContent.UID},
	); err != nil {
		return err
	}

	return err
}
