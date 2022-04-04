package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/models"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/repositories"
	customMongo "github.com/misterabdul/goblog-server/pkg/mongo"
)

// Get single post
func (service *Service) GetPost(filter interface{}) (
	post *models.PostModel,
	err error,
) {
	return repositories.GetPost(
		service.ctx,
		service.dbConn,
		filter)
}

// Get single post with its content
func (service *Service) GetPostWithContent(filter interface{}) (
	post *models.PostModel,
	content *models.PostContentModel,
	err error,
) {
	if post, err = repositories.GetPost(
		service.ctx, service.dbConn, filter,
	); err != nil {
		return nil, nil, err
	}
	if post == nil {
		return nil, nil, nil
	}
	if content, err = repositories.GetPostContent(
		service.ctx, service.dbConn, bson.M{
			"_id": bson.M{"$eq": post.UID}},
	); err != nil {
		return post, nil, err
	}

	return post, content, nil
}

// Get multiple posts
func (service *Service) GetPosts(filter interface{}) (
	posts []*models.PostModel,
	err error,
) {

	return repositories.GetPosts(
		service.ctx,
		service.dbConn,
		filter,
		internalGin.GetFindOptionsPost(service.c))
}

// Create new post with its content
func (service *Service) CreatePost(
	post *models.PostModel,
	content *models.PostContentModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	post.UID = primitive.NewObjectID()
	post.CreatedAt = now
	post.UpdatedAt = now
	post.DeletedAt = nil
	content.UID = post.UID

	return customMongo.Transaction(service.ctx, service.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = repositories.SavePost(
				sCtx, dbConn, post,
			); sErr != nil {
				return sErr
			}
			if sErr = repositories.SavePostContent(
				sCtx, dbConn, content,
			); sErr != nil {
				return sErr
			}

			return nil
		})

}

// Mark the post published
func (service *Service) PublishPost(
	post *models.PostModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	post.PublishedAt = now

	return repositories.UpdatePost(
		service.ctx,
		service.dbConn,
		post)
}

// Remove published mark from the post
func (service *Service) DepublishPost(
	post *models.PostModel,
) (err error) {
	post.PublishedAt = nil

	return repositories.UpdatePost(
		service.ctx,
		service.dbConn,
		post)
}

// Update post
func (service *Service) UpdatePost(
	post *models.PostModel,
	content *models.PostContentModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	post.UpdatedAt = now

	return customMongo.Transaction(service.ctx, service.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = repositories.UpdatePost(
				sCtx, dbConn, post,
			); sErr != nil {
				return sErr
			}
			if sErr = repositories.UpdatePostContent(
				sCtx, dbConn, content,
			); sErr != nil {
				return sErr
			}

			return nil
		})
}

// Delete post to trash
func (service *Service) TrashPost(
	post *models.PostModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	post.DeletedAt = now

	return repositories.UpdatePost(
		service.ctx,
		service.dbConn,
		post)
}

// Restore post from trash
func (service *Service) DetrashPost(
	post *models.PostModel,
) (err error) {
	post.DeletedAt = nil

	return repositories.UpdatePost(
		service.ctx,
		service.dbConn,
		post)
}

// Permanently delete post
func (service *Service) DeletePost(
	post *models.PostModel,
	content *models.PostContentModel,
) (err error) {

	return customMongo.Transaction(service.ctx, service.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = repositories.DeletePost(
				sCtx, dbConn, post,
			); sErr != nil {
				return sErr
			}
			if sErr = repositories.DeletePostContent(
				sCtx, dbConn, content,
			); sErr != nil {
				return sErr
			}

			return nil
		})
}
