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

// Get single post
func GetPost(ctx context.Context, dbConn *mongo.Database, filter interface{}) (*models.PostModel, error) {
	var post models.PostModel
	if err := getPostCollection(dbConn).FindOne(ctx, filter).Decode(&post); err != nil {
		return nil, err
	}

	return &post, nil
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
func CreatePost(ctx context.Context, dbConn *mongo.Database, post *models.PostModel) error {
	now := primitive.NewDateTimeFromTime(time.Now())

	post.UID = primitive.NewObjectID()
	post.CreatedAt = now
	post.UpdatedAt = now
	post.DeletedAt = nil

	insRes, err := getPostCollection(dbConn).InsertOne(ctx, post)
	if err != nil {
		return err
	}
	insertedID, ok := insRes.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("unable to assert inserted uid")
	}
	if post.UID != insertedID {
		return errors.New("inserted uid is not same with database")
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
func UpdatePost(ctx context.Context, dbConn *mongo.Database, post *models.PostModel) error {
	now := primitive.NewDateTimeFromTime(time.Now())

	post.UpdatedAt = now

	_, err := getPostCollection(dbConn).UpdateByID(ctx, post.UID, bson.M{"$set": post})

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
func DeletePost(ctx context.Context, dbConn *mongo.Database, post *models.PostModel) error {
	_, err := getPostCollection(dbConn).DeleteOne(ctx, bson.M{"_id": post.UID})

	return err
}
