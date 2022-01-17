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

// Get single post with its content
func GetPostWithContent(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
) (
	post *models.PostModel,
	postContent *models.PostContentModel,
	err error,
) {
	var (
		_post        models.PostModel
		_postContent models.PostContentModel
	)

	if err = getPostCollection(dbConn).FindOne(
		ctx, filter,
	).Decode(&_post); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	if err = getPostContentCollection(dbConn).FindOne(
		ctx, bson.M{"_id": _post.UID},
	).Decode(&_postContent); err != nil {
		if err == mongo.ErrNoDocuments {
			return &_post, nil, nil
		}
		return nil, nil, err
	}

	return &_post, &_postContent, err
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

// Save new post with its content
func SavePostWithContent(
	ctx context.Context,
	dbConn *mongo.Database,
	post *models.PostModel,
	postContent *models.PostContentModel,
) (err error) {
	var (
		session    mongo.Session
		insRes     *mongo.InsertOneResult
		insertedID interface{}
		ok         bool
	)

	if session, err = dbConn.Client().StartSession(); err != nil {
		return err
	}
	defer session.EndSession(ctx)
	if err = mongo.WithSession(ctx, session, func(sCtx mongo.SessionContext) (sErr error) {
		if sErr = sCtx.StartTransaction(); sErr != nil {
			return sErr
		}

		if insRes, sErr = getPostCollection(dbConn).InsertOne(
			sCtx, post,
		); sErr != nil {
			return sErr
		}
		if insertedID, ok = insRes.InsertedID.(primitive.ObjectID); !ok {
			return errors.New("unable to assert inserted uid")
		}
		if post.UID != insertedID {
			return errors.New("inserted uid is not same with database")
		}
		if insRes, sErr = getPostContentCollection(dbConn).InsertOne(
			sCtx, postContent,
		); sErr != nil {
			return sErr
		}
		if insertedID, ok = insRes.InsertedID.(primitive.ObjectID); !ok {
			return errors.New("unable to assert inserted uid")
		}
		if postContent.UID != insertedID {
			return errors.New("inserted uid is not same with database")
		}
		if sErr = session.CommitTransaction(sCtx); sErr != nil {
			return sErr
		}

		return nil
	}); err != nil {
		return err
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

// Update post with its content
func UpdatePostWithContent(
	ctx context.Context,
	dbConn *mongo.Database,
	post *models.PostModel,
	postContent *models.PostContentModel,
) (err error) {
	var session mongo.Session

	if post.UID != postContent.UID {
		return errors.New("post id not same as post content id")
	}
	if session, err = dbConn.Client().StartSession(); err != nil {
		return err
	}
	defer session.EndSession(ctx)
	if err = mongo.WithSession(ctx, session, func(sCtx mongo.SessionContext) (sErr error) {
		if sErr = sCtx.StartTransaction(); sErr != nil {
			return sErr
		}
		if _, sErr = getPostCollection(dbConn).UpdateByID(
			sCtx, post.UID, bson.M{"$set": post},
		); sErr != nil {
			return sErr
		}
		if _, sErr = getPostContentCollection(dbConn).UpdateByID(
			sCtx, postContent.UID, bson.M{"$set": postContent},
		); sErr != nil {
			return sErr
		}
		if sErr = session.CommitTransaction(sCtx); sErr != nil {
			return sErr
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// Update post
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

// Delete post with its content
func DeletePostWithContent(
	ctx context.Context,
	dbConn *mongo.Database,
	post *models.PostModel,
	postContent *models.PostContentModel,
) (err error) {
	var session mongo.Session

	if post.UID != postContent.UID {
		return errors.New("post id not same as post content id")
	}
	if session, err = dbConn.Client().StartSession(); err != nil {
		return err
	}
	defer session.EndSession(ctx)
	if err = mongo.WithSession(ctx, session, func(sCtx mongo.SessionContext) (sErr error) {
		if sErr = sCtx.StartTransaction(); sErr != nil {
			return sErr
		}

		if _, sErr = getPostCollection(dbConn).DeleteOne(
			sCtx, bson.M{"_id": post.UID},
		); sErr != nil {
			return sErr
		}

		if _, sErr = getPostContentCollection(dbConn).DeleteOne(
			sCtx, bson.M{"_id": postContent.UID},
		); sErr != nil {
			return sErr
		}

		if sErr = session.CommitTransaction(sCtx); sErr != nil {
			return sErr
		}

		return nil
	}); err != nil {
		return err
	}

	return err
}
