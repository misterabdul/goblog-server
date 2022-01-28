package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/models"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/repositories"
	customMongo "github.com/misterabdul/goblog-server/pkg/mongo"
)

// Get single comment
func (service *Service) GetComment(
	filter interface{},
) (comment *models.CommentModel, err error) {

	return repositories.GetComment(
		service.ctx,
		service.dbConn,
		filter)
}

// Get multiple comments
func (service *Service) GetComments(
	filter interface{},
) (comments []*models.CommentModel, err error) {

	return repositories.GetComments(
		service.ctx,
		service.dbConn,
		filter,
		internalGin.GetFindOptions(service.c))
}

// Create new comment
func (service *Service) CreateComment(
	comment *models.CommentModel,
	post *models.PostModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	comment.UID = primitive.NewObjectID()
	comment.CreatedAt = now
	comment.DeletedAt = nil
	post.CommentCount++

	return customMongo.Transaction(service.ctx, service.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = repositories.SaveComment(sCtx, dbConn, comment); err != nil {
				return err
			}
			if err = repositories.UpdatePost(sCtx, dbConn, post); err != nil {
				return err
			}

			return nil
		})
}

// Create new comment reply
func (service *Service) CreateCommentReply(
	reply *models.CommentModel,
	comment *models.CommentModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	reply.UID = primitive.NewObjectID()
	reply.CreatedAt = now
	reply.DeletedAt = nil
	comment.ReplyCount++

	return customMongo.Transaction(service.ctx, service.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = repositories.SaveComment(sCtx, dbConn, reply); err != nil {
				return nil
			}
			if err = repositories.SaveComment(sCtx, dbConn, comment); err != nil {
				return nil
			}

			return nil
		})
}

// Delete comment to trash
func (service *Service) TrashComment(
	comment *models.CommentModel,
	post *models.PostModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	comment.DeletedAt = now
	post.CommentCount--

	return customMongo.Transaction(service.ctx, service.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = repositories.UpdateComment(sCtx, dbConn, comment); err != nil {
				return nil
			}
			if err = repositories.UpdatePost(sCtx, dbConn, post); err != nil {
				return err
			}

			return nil
		})
}

// Delete comment reply to trash
func (service *Service) TrashCommentReply(
	reply *models.CommentModel,
	comment *models.CommentModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	reply.DeletedAt = now
	comment.ReplyCount--

	return customMongo.Transaction(service.ctx, service.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = repositories.UpdateComment(sCtx, dbConn, reply); err != nil {
				return nil
			}
			if err = repositories.UpdateComment(sCtx, dbConn, comment); err != nil {
				return err
			}

			return nil
		})
}

// Restore comment from trash
func (service *Service) DetrashComment(
	comment *models.CommentModel,
	post *models.PostModel,
) (err error) {
	comment.DeletedAt = nil
	post.CommentCount++

	return customMongo.Transaction(service.ctx, service.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = repositories.UpdateComment(sCtx, dbConn, comment); err != nil {
				return nil
			}
			if err = repositories.UpdatePost(sCtx, dbConn, post); err != nil {
				return err
			}

			return nil
		})
}

// Restore comment reply from trash
func (service *Service) DetrashCommentReply(
	reply *models.CommentModel,
	comment *models.CommentModel,
) (err error) {
	reply.DeletedAt = nil
	comment.ReplyCount++

	return customMongo.Transaction(service.ctx, service.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = repositories.UpdateComment(sCtx, dbConn, reply); err != nil {
				return nil
			}
			if err = repositories.UpdateComment(sCtx, dbConn, comment); err != nil {
				return err
			}

			return nil
		})
}

// Permanently delete comment
func (service *Service) DeleteComment(
	comment *models.CommentModel,
	post *models.PostModel,
) (err error) {
	post.CommentCount--

	return customMongo.Transaction(service.ctx, service.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = repositories.DeleteComment(sCtx, dbConn, comment); err != nil {
				return nil
			}
			if err = repositories.UpdatePost(sCtx, dbConn, post); err != nil {
				return err
			}

			return nil
		})
}

// Permanently delete comment
func (service *Service) DeleteCommentReply(
	reply *models.CommentModel,
	comment *models.CommentModel,
) (err error) {
	comment.ReplyCount--

	return customMongo.Transaction(service.ctx, service.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = repositories.DeleteComment(sCtx, dbConn, reply); err != nil {
				return nil
			}
			if err = repositories.UpdateComment(sCtx, dbConn, comment); err != nil {
				return err
			}

			return nil
		})
}
