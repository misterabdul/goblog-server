package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/database/repositories"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	customMongo "github.com/misterabdul/goblog-server/pkg/mongo"
)

type PostService struct {
	c                 *gin.Context
	ctx               context.Context
	dbConn            *mongo.Database
	repository        *repositories.PostRepository
	contentRepository *repositories.PostContentRepository
}

func NewPostService(
	c *gin.Context,
	ctx context.Context,
	dbConn *mongo.Database,
) *PostService {

	return &PostService{
		c:                 c,
		ctx:               ctx,
		dbConn:            dbConn,
		repository:        repositories.NewPostRepository(dbConn),
		contentRepository: repositories.NewPostContentRepository(dbConn)}
}

// Get single post
func (s *PostService) GetPost(filter interface{}) (
	post *models.PostModel,
	err error,
) {
	return s.repository.ReadOne(
		s.ctx, filter)
}

// Get single post with its content
func (s *PostService) GetPostWithContent(filter interface{}) (
	post *models.PostModel,
	content *models.PostContentModel,
	err error,
) {
	if post, err = s.repository.ReadOne(
		s.ctx, filter,
	); err != nil {
		return nil, nil, err
	}
	if post == nil {
		return nil, nil, nil
	}
	if content, err = s.contentRepository.ReadOne(
		s.ctx, bson.M{
			"_id": bson.M{"$eq": post.UID}},
	); err != nil {
		return post, nil, err
	}

	return post, content, nil
}

// Get multiple posts
func (s *PostService) GetPosts(filter interface{}) (
	posts []*models.PostModel,
	err error,
) {

	return s.repository.ReadMany(
		s.ctx, filter,
		internalGin.GetFindOptionsPost(s.c))
}

// Get total posts count
func (s *PostService) GetPostCount(filter interface{}) (
	count int64, err error,
) {

	return s.repository.Count(
		s.ctx, filter,
		internalGin.GetCountOptions(s.c))
}

// Create new post with its content
func (s *PostService) CreatePost(
	post *models.PostModel,
	content *models.PostContentModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	post.UID = primitive.NewObjectID()
	post.CreatedAt = now
	post.UpdatedAt = now
	post.DeletedAt = nil
	content.UID = post.UID

	return customMongo.Transaction(s.ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = s.repository.Save(
				sCtx, post,
			); sErr != nil {
				return sErr
			}
			if sErr = s.contentRepository.Save(
				sCtx, content,
			); sErr != nil {
				return sErr
			}

			return nil
		})

}

// Mark the post published
func (s *PostService) PublishPost(
	post *models.PostModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	post.PublishedAt = now

	return s.repository.Update(
		s.ctx, post)
}

// Remove published mark from the post
func (s *PostService) DepublishPost(
	post *models.PostModel,
) (err error) {
	post.PublishedAt = nil

	return s.repository.Update(
		s.ctx, post)
}

// Update post
func (s *PostService) UpdatePost(
	post *models.PostModel,
	content *models.PostContentModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	post.UpdatedAt = now

	return customMongo.Transaction(s.ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = s.repository.Update(
				sCtx, post,
			); sErr != nil {
				return sErr
			}
			if sErr = s.contentRepository.Update(
				sCtx, content,
			); sErr != nil {
				return sErr
			}

			return nil
		})
}

// Delete post to trash
func (s *PostService) TrashPost(
	post *models.PostModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	post.DeletedAt = now

	return s.repository.Update(
		s.ctx, post)
}

// Restore post from trash
func (s *PostService) DetrashPost(
	post *models.PostModel,
) (err error) {
	post.DeletedAt = nil

	return s.repository.Update(
		s.ctx, post)
}

// Permanently delete post
func (s *PostService) DeletePost(
	post *models.PostModel,
	content *models.PostContentModel,
) (err error) {

	return customMongo.Transaction(s.ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if sErr = s.repository.Delete(
				sCtx, post,
			); sErr != nil {
				return sErr
			}
			if sErr = s.contentRepository.Delete(
				sCtx, content,
			); sErr != nil {
				return sErr
			}

			return nil
		})
}
