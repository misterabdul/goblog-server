package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/database/repositories"
	customMongo "github.com/misterabdul/goblog-server/pkg/mongo"
)

type comment struct {
	dbConn *mongo.Database
}

func newCommentService(
	dbConn *mongo.Database,
) (service *comment) {

	return &comment{dbConn: dbConn}
}

// Get single comment
func (s *comment) GetOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (comment *models.CommentModel, err error) {

	return repositories.ReadOneComment(
		s.dbConn, ctx, filter)
}

// Get multiple comments
func (s *comment) GetMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (comments []*models.CommentModel, err error) {

	return repositories.ReadManyComments(
		s.dbConn, ctx, filter, opts...)
}

// Get total comments count
func (s *comment) Count(
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (count int64, err error) {

	return repositories.CountComments(
		s.dbConn, ctx, filter, opts...)
}

// Create new comment
func (s *comment) SaveOne(
	ctx context.Context,
	comment *models.CommentModel,
	post *models.PostModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	comment.UID = primitive.NewObjectID()
	comment.CreatedAt = now
	comment.DeletedAt = nil
	post.CommentCount++

	return customMongo.Transaction(ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = repositories.SaveOneComment(dbConn, sCtx, comment, opts...); err != nil {
				return err
			}
			if err = repositories.UpdateOnePost(dbConn, sCtx, post); err != nil {
				return err
			}

			return nil
		})
}

// Create new comment reply
func (s *comment) SaveOneReply(
	ctx context.Context,
	reply *models.CommentModel,
	comment *models.CommentModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	reply.UID = primitive.NewObjectID()
	reply.CreatedAt = now
	reply.DeletedAt = nil
	comment.ReplyCount++

	return customMongo.Transaction(ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = repositories.SaveOneComment(dbConn, sCtx, reply, opts...); err != nil {
				return nil
			}
			if err = repositories.UpdateOneComment(dbConn, sCtx, comment); err != nil {
				return nil
			}

			return nil
		})
}

// Delete comment to trash
func (s *comment) TrashOne(
	ctx context.Context,
	comment *models.CommentModel,
	post *models.PostModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	comment.DeletedAt = now
	post.CommentCount--

	return customMongo.Transaction(ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = repositories.UpdateOneComment(dbConn, sCtx, comment, opts...); err != nil {
				return nil
			}
			if err = repositories.UpdateOnePost(dbConn, sCtx, post); err != nil {
				return err
			}

			return nil
		})
}

// Delete comment reply to trash
func (s *comment) TrashOneReply(
	ctx context.Context,
	reply *models.CommentModel,
	comment *models.CommentModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	reply.DeletedAt = now
	comment.ReplyCount--

	return customMongo.Transaction(ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = repositories.UpdateOneComment(dbConn, sCtx, reply, opts...); err != nil {
				return nil
			}
			if err = repositories.UpdateOneComment(dbConn, sCtx, comment); err != nil {
				return err
			}

			return nil
		})
}

// Restore comment from trash
func (s *comment) RestoreOne(
	ctx context.Context,
	comment *models.CommentModel,
	post *models.PostModel,
	opts ...*options.UpdateOptions,
) (err error) {
	comment.DeletedAt = nil
	post.CommentCount++

	return customMongo.Transaction(ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = repositories.UpdateOneComment(dbConn, sCtx, comment, opts...); err != nil {
				return nil
			}
			if err = repositories.UpdateOnePost(dbConn, sCtx, post); err != nil {
				return err
			}

			return nil
		})
}

// Restore comment reply from trash
func (s *comment) RestoreOneReply(
	ctx context.Context,
	reply *models.CommentModel,
	comment *models.CommentModel,
	opts ...*options.UpdateOptions,
) (err error) {
	reply.DeletedAt = nil
	comment.ReplyCount++

	return customMongo.Transaction(ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = repositories.UpdateOneComment(dbConn, sCtx, reply, opts...); err != nil {
				return nil
			}
			if err = repositories.UpdateOneComment(dbConn, sCtx, comment); err != nil {
				return err
			}

			return nil
		})
}

// Permanently delete comment
func (s *comment) DeleteOne(
	ctx context.Context,
	comment *models.CommentModel,
	post *models.PostModel,
	opts ...*options.DeleteOptions,
) (err error) {
	post.CommentCount--

	return customMongo.Transaction(ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = repositories.DeleteOneComment(dbConn, sCtx, comment, opts...); err != nil {
				return nil
			}
			if err = repositories.UpdateOnePost(dbConn, sCtx, post); err != nil {
				return err
			}

			return nil
		})
}

// Permanently delete comment
func (s *comment) DeleteOneReply(
	ctx context.Context,
	reply *models.CommentModel,
	comment *models.CommentModel,
	opts ...*options.DeleteOptions,
) (err error) {
	comment.ReplyCount--

	return customMongo.Transaction(ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = repositories.DeleteOneComment(dbConn, sCtx, reply, opts...); err != nil {
				return nil
			}
			if err = repositories.UpdateOneComment(dbConn, sCtx, comment); err != nil {
				return err
			}

			return nil
		})
}
