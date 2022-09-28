package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/database/repositories"
	customMongo "github.com/misterabdul/goblog-server/pkg/mongo"
)

type post struct {
	dbConn *mongo.Database
}

func newPostService(
	dbConn *mongo.Database,
) (service *post) {

	return &post{dbConn: dbConn}
}

// Get single post
func (s *post) GetOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (post *models.PostModel, err error) {

	return repositories.ReadOnePost(
		s.dbConn, ctx, filter, opts...)
}

// Get single post with its content
func (s *post) GetOneWithContent(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (post *models.PostModel, content *models.PostContentModel, err error) {
	if post, err = repositories.ReadOnePost(
		s.dbConn, ctx, filter, opts...,
	); err != nil {
		return nil, nil, err
	}
	if post == nil {
		return nil, nil, nil
	}
	if content, err = repositories.ReadOnePostContent(
		s.dbConn, ctx, bson.M{"_id": bson.M{"$eq": post.UID}}, opts...,
	); err != nil {
		return post, nil, err
	}

	return post, content, nil
}

// Get multiple posts
func (s *post) GetMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (posts []*models.PostModel, err error) {

	return repositories.ReadManyPosts(
		s.dbConn, ctx, filter, opts...)
}

// Get total posts count
func (s *post) Count(
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {

	return repositories.CountPosts(
		s.dbConn, ctx, filter, opts...)
}

// Create new post with its content
func (s *post) SaveOneWithContent(
	ctx context.Context,
	post *models.PostModel,
	content *models.PostContentModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	post.UID = primitive.NewObjectID()
	post.CreatedAt = now
	post.UpdatedAt = now
	post.DeletedAt = nil
	content.UID = post.UID

	return customMongo.Transaction(ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = repositories.SaveOnePost(
				dbConn, sCtx, post, opts...,
			); sErr != nil {
				return sErr
			}
			if sErr = repositories.SaveOnePostContent(
				dbConn, sCtx, content, opts...,
			); sErr != nil {
				return sErr
			}

			return nil
		})

}

// Mark the post published
func (s *post) PublishOne(
	ctx context.Context,
	post *models.PostModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	post.PublishedAt = now

	return repositories.UpdateOnePost(
		s.dbConn, ctx, post, opts...)
}

// Remove published mark from the post
func (s *post) DepublishOne(
	ctx context.Context,
	post *models.PostModel,
	opts ...*options.UpdateOptions,
) (err error) {
	post.PublishedAt = nil

	return repositories.UpdateOnePost(
		s.dbConn, ctx, post, opts...)
}

// Update post
func (s *post) UpdateOneWithContent(
	ctx context.Context,
	post *models.PostModel,
	content *models.PostContentModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	post.UpdatedAt = now

	return customMongo.Transaction(ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = repositories.UpdateOnePost(
				dbConn, sCtx, post,
			); sErr != nil {
				return sErr
			}
			if sErr = repositories.UpdateOnePostContent(
				dbConn, sCtx, content,
			); sErr != nil {
				return sErr
			}

			return nil
		})
}

// Update post's author
func (s *post) UpdateManyAuthor(
	ctx context.Context,
	author *models.UserModel,
	opts ...*options.UpdateOptions,
) (err error) {
	return customMongo.Transaction(ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = repositories.UpdateManyPostAuthor(
				dbConn, sCtx, author,
			); sErr != nil {
				return sErr
			}

			return nil
		})
}

// Delete post to trash
func (s *post) TrashOne(
	ctx context.Context,
	post *models.PostModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	post.DeletedAt = now

	return repositories.UpdateOnePost(
		s.dbConn, ctx, post)
}

// Restore post from trash
func (s *post) RestoreOne(
	ctx context.Context,
	post *models.PostModel,
	opts ...*options.UpdateOptions,
) (err error) {
	post.DeletedAt = nil

	return repositories.UpdateOnePost(
		s.dbConn, ctx, post, opts...)
}

// Permanently delete post
func (s *post) DeleteOneWithContent(
	ctx context.Context,
	post *models.PostModel,
	content *models.PostContentModel,
	opts ...*options.DeleteOptions,
) (err error) {

	return customMongo.Transaction(ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = repositories.DeleteOnePost(
				dbConn, sCtx, post,
			); sErr != nil {
				return sErr
			}
			if sErr = repositories.DeleteOnePostContent(
				dbConn, sCtx, content,
			); sErr != nil {
				return sErr
			}

			return nil
		})
}
