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

func getPostCollection(
	dbConn *mongo.Database,
) (postCollection *mongo.Collection) {
	return dbConn.Collection("posts")
}

func getPostContentCollection(dbConn *mongo.Database,
) (postContentCollection *mongo.Collection) {
	return dbConn.Collection("postContents")
}

// Get single post
func GetPost(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.FindOneOptions,
) (post *models.PostModel, err error) {
	var _post models.PostModel

	if err = getPostCollection(dbConn).FindOne(
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
func GetPostContent(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.FindOneOptions,
) (
	postContent *models.PostContentModel,
	err error,
) {
	var _postContent models.PostContentModel

	if err = getPostContentCollection(dbConn).FindOne(
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
func GetPosts(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.FindOptions,
) (posts []*models.PostModel, err error) {
	var (
		post   *models.PostModel
		cursor *mongo.Cursor
	)

	if cursor, err = getPostCollection(dbConn).Find(
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
func CountPosts(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {

	return getPostCollection(dbConn).CountDocuments(
		ctx, filter, opts...,
	)
}

// Save new post
func SavePost(
	ctx context.Context,
	dbConn *mongo.Database,
	post *models.PostModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID interface{}
		ok         bool
	)

	if insRes, err = getPostCollection(dbConn).InsertOne(
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
func SavePostContent(
	ctx context.Context,
	dbConn *mongo.Database,
	postContent *models.PostContentModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID interface{}
		ok         bool
	)

	if insRes, err = getPostContentCollection(dbConn).InsertOne(
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
func UpdatePost(
	ctx context.Context,
	dbConn *mongo.Database,
	post *models.PostModel,
) (err error) {
	if _, err = getPostCollection(dbConn).UpdateByID(
		ctx, post.UID, bson.M{"$set": post},
	); err != nil {
		return err
	}

	return nil
}

// Update post content
func UpdatePostContent(
	ctx context.Context,
	dbConn *mongo.Database,
	postContent *models.PostContentModel,
) (err error) {
	if _, err = getPostContentCollection(dbConn).UpdateByID(
		ctx, postContent.UID, bson.M{"$set": postContent},
	); err != nil {
		return err
	}

	return nil
}

// Delete post
func DeletePost(
	ctx context.Context,
	dbConn *mongo.Database,
	post *models.PostModel,
) (err error) {
	if _, err = getPostCollection(dbConn).DeleteOne(
		ctx, bson.M{"_id": post.UID},
	); err != nil {
		return err
	}

	return nil
}

// Delete post content
func DeletePostContent(
	ctx context.Context,
	dbConn *mongo.Database,
	postContent *models.PostContentModel,
) (err error) {
	if _, err = getPostContentCollection(dbConn).DeleteOne(
		ctx, bson.M{"_id": postContent.UID},
	); err != nil {
		return err
	}

	return err
}
