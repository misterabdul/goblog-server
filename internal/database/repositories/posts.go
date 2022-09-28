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

const (
	postCollection        = "posts"
	postContentCollection = "postContents"
)

// Get single post
func ReadOnePost(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (post *models.PostModel, err error) {
	var (
		collection = dbConn.Collection(postCollection)
		_post      models.PostModel
	)

	if err = collection.FindOne(ctx, filter, opts...).Decode(&_post); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_post, nil
}

// Get single post content
func ReadOnePostContent(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (
	postContent *models.PostContentModel,
	err error,
) {
	var (
		collection   = dbConn.Collection(postContentCollection)
		_postContent models.PostContentModel
	)

	if err = collection.FindOne(ctx, filter, opts...).Decode(&_postContent); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_postContent, nil
}

// Get multiple posts
func ReadManyPosts(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (posts []*models.PostModel, err error) {
	var (
		collection = dbConn.Collection(postCollection)
		post       *models.PostModel
		cursor     *mongo.Cursor
	)

	if cursor, err = collection.Find(ctx, filter, opts...); err != nil {
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
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {
	var collection = dbConn.Collection(postCollection)

	return collection.CountDocuments(
		ctx, filter, opts...)
}

// Save new post
func SaveOnePost(
	dbConn *mongo.Database,
	ctx context.Context,
	post *models.PostModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var (
		collection = dbConn.Collection(postCollection)
		insRes     *mongo.InsertOneResult
		insertedID interface{}
		ok         bool
	)

	if insRes, err = collection.InsertOne(ctx, post, opts...); err != nil {
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
func SaveOnePostContent(
	dbConn *mongo.Database,
	ctx context.Context,
	postContent *models.PostContentModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var (
		collection = dbConn.Collection(postContentCollection)
		insRes     *mongo.InsertOneResult
		insertedID interface{}
		ok         bool
	)

	if insRes, err = collection.InsertOne(ctx, postContent, opts...); err != nil {
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
func UpdateOnePost(
	dbConn *mongo.Database,
	ctx context.Context,
	post *models.PostModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var collection = dbConn.Collection(postCollection)

	_, err = collection.UpdateOne(
		ctx, bson.M{"_id": post.UID}, bson.M{"$set": post}, opts...)

	return err
}

// Bulk update post's author
func UpdateManyPostAuthor(
	dbConn *mongo.Database,
	ctx context.Context,
	user *models.UserModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var collection = dbConn.Collection(postCollection)

	_, err = collection.UpdateMany(ctx,
		bson.M{"author._id": bson.M{"$eq": user.UID}},
		bson.M{"$set": bson.M{"author": user.ToCommonModel()}}, opts...)

	return err
}

// Update post content
func UpdateOnePostContent(
	dbConn *mongo.Database,
	ctx context.Context,
	postContent *models.PostContentModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var collection = dbConn.Collection(postContentCollection)

	_, err = collection.UpdateOne(
		ctx, bson.M{"_id": postContent.UID}, bson.M{"$set": postContent}, opts...)

	return err
}

// Delete post
func DeleteOnePost(
	dbConn *mongo.Database,
	ctx context.Context,
	post *models.PostModel,
	opts ...*options.DeleteOptions,
) (err error) {
	var collection = dbConn.Collection(postContentCollection)

	_, err = collection.DeleteOne(
		ctx, bson.M{"_id": post.UID}, opts...)

	return err
}

// Delete post content
func DeleteOnePostContent(
	dbConn *mongo.Database,
	ctx context.Context,
	postContent *models.PostContentModel,
	opts ...*options.DeleteOptions,
) (err error) {
	var collection = dbConn.Collection(postContentCollection)

	_, err = collection.DeleteOne(
		ctx, bson.M{"_id": postContent.UID}, opts...)

	return err
}
