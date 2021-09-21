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

func getPostCollection(dbConn *mongo.Database) *mongo.Collection {
	return dbConn.Collection("posts")
}

func getPostContentCollection(dbConn *mongo.Database) *mongo.Collection {
	return dbConn.Collection("postContents")
}

// Get single post
func GetPost(ctx context.Context, dbConn *mongo.Database, filter interface{}) (*models.PostModel, error) {
	var (
		post models.PostModel
		err  error
	)

	if err = getPostCollection(dbConn).FindOne(ctx, filter).Decode(&post); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &post, nil
}

// Get single post with its content
func GetPostWithContent(ctx context.Context, dbConn *mongo.Database, filter interface{}) (*models.PostModel, *models.PostContentModel, error) {
	var (
		post        models.PostModel
		postContent models.PostContentModel
		err         error
	)
	if err = getPostCollection(dbConn).FindOne(ctx, filter).Decode(&post); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	if err = getPostContentCollection(dbConn).FindOne(ctx, bson.M{"_id": post.UID}).Decode(&postContent); err != nil {
		if err == mongo.ErrNoDocuments {
			return &post, nil, err
		}
	}

	return &post, &postContent, err
}

// Get multiple posts
func GetPosts(ctx context.Context, dbConn *mongo.Database, filter interface{}, show int, order string, asc bool) ([]*models.PostModel, error) {
	cursor, err := getPostCollection(dbConn).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*models.PostModel
	for cursor.Next(ctx) {
		var post models.PostModel
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

// Create new post
func CreatePost(ctx context.Context, dbConn *mongo.Database, post *models.PostModel, postContent *models.PostContentModel) error {
	var (
		now        = primitive.NewDateTimeFromTime(time.Now())
		session    mongo.Session
		insRes     *mongo.InsertOneResult
		insertedID interface{}
		ok         bool
		err        error
	)

	if session, err = dbConn.Client().StartSession(); err != nil {
		return err
	}
	defer session.EndSession(ctx)

	if err = mongo.WithSession(ctx, session, func(sctx mongo.SessionContext) error {
		if err = sctx.StartTransaction(); err != nil {
			return err
		}

		post.UID = primitive.NewObjectID()
		post.CreatedAt = now
		post.UpdatedAt = now
		post.DeletedAt = nil
		if insRes, err = getPostCollection(dbConn).InsertOne(sctx, post); err != nil {
			return err
		}
		if insertedID, ok = insRes.InsertedID.(primitive.ObjectID); !ok {
			return errors.New("unable to assert inserted uid")
		}
		if post.UID != insertedID {
			return errors.New("inserted uid is not same with database")
		}

		postContent.UID = post.UID
		if insRes, err = getPostContentCollection(dbConn).InsertOne(sctx, postContent); err != nil {
			return err
		}
		if insertedID, ok = insRes.InsertedID.(primitive.ObjectID); !ok {
			return errors.New("unable to assert inserted uid")
		}
		if postContent.UID != insertedID {
			return errors.New("inserted uid is not same with database")
		}

		if err = session.CommitTransaction(sctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// Publish the post
func PublishPost(ctx context.Context, dbConn *mongo.Database, post *models.PostModel) error {
	now := primitive.NewDateTimeFromTime(time.Now())

	post.PublishedAt = now

	_, err := getPostCollection(dbConn).UpdateByID(ctx, post.UID, bson.M{"$set": post})

	return err
}

// Depublish the post
func DepublishPost(ctx context.Context, dbConn *mongo.Database, post *models.PostModel) error {
	post.PublishedAt = nil

	_, err := getPostCollection(dbConn).UpdateByID(ctx, post.UID, bson.M{"$set": post})

	return err
}

// Update post
func UpdatePost(ctx context.Context, dbConn *mongo.Database, post *models.PostModel, postContent *models.PostContentModel) error {
	var (
		now     = primitive.NewDateTimeFromTime(time.Now())
		session mongo.Session
		err     error
	)

	if post.UID != postContent.UID {
		return errors.New("post id not same as post content id")
	}
	if session, err = dbConn.Client().StartSession(); err != nil {
		return err
	}
	defer session.EndSession(ctx)

	if err = mongo.WithSession(ctx, session, func(sctx mongo.SessionContext) error {
		if err = sctx.StartTransaction(); err != nil {
			return err
		}

		post.UpdatedAt = now
		if _, err = getPostCollection(dbConn).UpdateByID(sctx, post.UID, bson.M{"$set": post}); err != nil {
			return err
		}

		if _, err = getPostContentCollection(dbConn).UpdateByID(sctx, postContent.UID, bson.M{"$set": postContent}); err != nil {
			return err
		}

		if err = session.CommitTransaction(sctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return err
}

// Mark the post trash
func TrashPost(ctx context.Context, dbConn *mongo.Database, post *models.PostModel) error {
	now := primitive.NewDateTimeFromTime(time.Now())

	post.DeletedAt = now

	_, err := getPostCollection(dbConn).UpdateByID(ctx, post.UID, bson.M{"$set": post})

	return err
}

// Unmark the trash from post
func DetrashPost(ctx context.Context, dbConn *mongo.Database, post *models.PostModel) error {
	post.DeletedAt = nil

	_, err := getPostCollection(dbConn).UpdateByID(ctx, post.UID, bson.M{"$set": post})

	return err
}

// Permanently delete post
func DeletePost(ctx context.Context, dbConn *mongo.Database, post *models.PostModel, postContent *models.PostContentModel) error {
	var (
		session mongo.Session
		err     error
	)

	if post.UID != postContent.UID {
		return errors.New("post id not same as post content id")
	}
	if session, err = dbConn.Client().StartSession(); err != nil {
		return err
	}
	defer session.EndSession(ctx)

	if err = mongo.WithSession(ctx, session, func(sctx mongo.SessionContext) error {
		if err = sctx.StartTransaction(); err != nil {
			return err
		}

		if _, err = getPostCollection(dbConn).DeleteOne(sctx, bson.M{"_id": post.UID}); err != nil {
			return err
		}

		if _, err = getPostContentCollection(dbConn).DeleteOne(sctx, bson.M{"_id": postContent.UID}); err != nil {
			return err
		}

		if err = session.CommitTransaction(sctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return err
}
