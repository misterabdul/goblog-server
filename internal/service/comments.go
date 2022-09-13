package service

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/database/repositories"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	customMongo "github.com/misterabdul/goblog-server/pkg/mongo"
)

type CommentService struct {
	c              *gin.Context
	ctx            context.Context
	dbConn         *mongo.Database
	repository     *repositories.CommentRepository
	postRepository *repositories.PostRepository
}

func NewCommentService(
	c *gin.Context,
	ctx context.Context,
	dbConn *mongo.Database,
) *CommentService {

	return &CommentService{
		c:              c,
		ctx:            ctx,
		dbConn:         dbConn,
		repository:     repositories.NewCommentRepository(dbConn),
		postRepository: repositories.NewPostRepository(dbConn)}
}

// Get single comment
func (s *CommentService) GetComment(
	filter interface{},
) (comment *models.CommentModel, err error) {

	return s.repository.ReadOne(
		s.ctx, filter)
}

// Get multiple comments
func (s *CommentService) GetComments(
	filter interface{},
	dateDesc bool,
) (comments []*models.CommentModel, err error) {
	var _options *options.FindOptions

	if dateDesc {
		_options = &options.FindOptions{
			Sort: bson.M{"createdat": 1},
		}
	} else {
		_options = internalGin.GetFindOptions(s.c)
	}

	return s.repository.ReadMany(
		s.ctx, filter, _options)
}

// Get total comments count
func (s *CommentService) GetCommentCount(filter interface{}) (
	count int64, err error,
) {

	return s.repository.Count(
		s.ctx, filter,
		internalGin.GetCountOptions(s.c))
}

// Create new comment
func (s *CommentService) CreateComment(
	comment *models.CommentModel,
	post *models.PostModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	comment.UID = primitive.NewObjectID()
	comment.CreatedAt = now
	comment.DeletedAt = nil
	post.CommentCount++

	return customMongo.Transaction(s.ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = s.repository.Save(sCtx, comment); err != nil {
				return err
			}
			if err = s.postRepository.Update(sCtx, post); err != nil {
				return err
			}

			return nil
		})
}

// Create new comment reply
func (s *CommentService) CreateCommentReply(
	reply *models.CommentModel,
	comment *models.CommentModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	reply.UID = primitive.NewObjectID()
	reply.CreatedAt = now
	reply.DeletedAt = nil
	comment.ReplyCount++

	return customMongo.Transaction(s.ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = s.repository.Save(sCtx, reply); err != nil {
				return nil
			}
			if err = s.repository.Update(sCtx, comment); err != nil {
				return nil
			}

			return nil
		})
}

// Delete comment to trash
func (s *CommentService) TrashComment(
	comment *models.CommentModel,
	post *models.PostModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	comment.DeletedAt = now
	post.CommentCount--

	return customMongo.Transaction(s.ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = s.repository.Update(sCtx, comment); err != nil {
				return nil
			}
			if err = s.postRepository.Update(sCtx, post); err != nil {
				return err
			}

			return nil
		})
}

// Delete comment reply to trash
func (s *CommentService) TrashCommentReply(
	reply *models.CommentModel,
	comment *models.CommentModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	reply.DeletedAt = now
	comment.ReplyCount--

	return customMongo.Transaction(s.ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = s.repository.Update(sCtx, reply); err != nil {
				return nil
			}
			if err = s.repository.Update(sCtx, comment); err != nil {
				return err
			}

			return nil
		})
}

// Restore comment from trash
func (s *CommentService) DetrashComment(
	comment *models.CommentModel,
	post *models.PostModel,
) (err error) {
	comment.DeletedAt = nil
	post.CommentCount++

	return customMongo.Transaction(s.ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = s.repository.Update(sCtx, comment); err != nil {
				return nil
			}
			if err = s.postRepository.Update(sCtx, post); err != nil {
				return err
			}

			return nil
		})
}

// Restore comment reply from trash
func (s *CommentService) DetrashCommentReply(
	reply *models.CommentModel,
	comment *models.CommentModel,
) (err error) {
	reply.DeletedAt = nil
	comment.ReplyCount++

	return customMongo.Transaction(s.ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = s.repository.Update(sCtx, reply); err != nil {
				return nil
			}
			if err = s.repository.Update(sCtx, comment); err != nil {
				return err
			}

			return nil
		})
}

// Permanently delete comment
func (s *CommentService) DeleteComment(
	comment *models.CommentModel,
	post *models.PostModel,
) (err error) {
	post.CommentCount--

	return customMongo.Transaction(s.ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = s.repository.Delete(sCtx, comment); err != nil {
				return nil
			}
			if err = s.postRepository.Update(sCtx, post); err != nil {
				return err
			}

			return nil
		})
}

// Permanently delete comment
func (s *CommentService) DeleteCommentReply(
	reply *models.CommentModel,
	comment *models.CommentModel,
) (err error) {
	comment.ReplyCount--

	return customMongo.Transaction(s.ctx, s.dbConn, false,
		func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
			if err = s.repository.Delete(sCtx, reply); err != nil {
				return nil
			}
			if err = s.repository.Update(sCtx, comment); err != nil {
				return err
			}

			return nil
		})
}
