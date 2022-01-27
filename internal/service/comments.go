package service

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/models"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/repositories"
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
	var (
		ctx     = service.ctx
		dbConn  = service.dbConn
		session mongo.Session
		now     = primitive.NewDateTimeFromTime(time.Now())
	)

	comment.UID = primitive.NewObjectID()
	comment.CreatedAt = now
	comment.DeletedAt = nil
	post.CommentCount++
	if session, err = dbConn.Client().StartSession(); err != nil {
		return err
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(sCtx mongo.SessionContext) (sErr error) {
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
	var (
		ctx     = service.ctx
		dbConn  = service.dbConn
		session mongo.Session
		now     = primitive.NewDateTimeFromTime(time.Now())
	)

	reply.UID = primitive.NewObjectID()
	reply.CreatedAt = now
	reply.DeletedAt = nil
	comment.ReplyCount++
	if session, err = dbConn.Client().StartSession(); err != nil {
		return err
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(sCtx mongo.SessionContext) (sErr error) {
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
	var (
		ctx     = service.ctx
		dbConn  = service.dbConn
		session mongo.Session
		now     = primitive.NewDateTimeFromTime(time.Now())
	)

	comment.DeletedAt = now
	post.CommentCount--
	if session, err = dbConn.Client().StartSession(); err != nil {
		return err
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(sCtx mongo.SessionContext) (sErr error) {
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
	var (
		ctx     = service.ctx
		dbConn  = service.dbConn
		session mongo.Session
		now     = primitive.NewDateTimeFromTime(time.Now())
	)

	reply.DeletedAt = now
	comment.ReplyCount--
	if session, err = dbConn.Client().StartSession(); err != nil {
		return err
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(sCtx mongo.SessionContext) (sErr error) {
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
	var (
		ctx     = service.ctx
		dbConn  = service.dbConn
		session mongo.Session
	)

	comment.DeletedAt = nil
	post.CommentCount++
	if session, err = dbConn.Client().StartSession(); err != nil {
		return err
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(sCtx mongo.SessionContext) (sErr error) {
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
	var (
		ctx     = service.ctx
		dbConn  = service.dbConn
		session mongo.Session
	)

	reply.DeletedAt = nil
	comment.ReplyCount++
	if session, err = dbConn.Client().StartSession(); err != nil {
		return err
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(sCtx mongo.SessionContext) (sErr error) {
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
	var (
		ctx     = service.ctx
		dbConn  = service.dbConn
		session mongo.Session
	)

	post.CommentCount--
	if session, err = dbConn.Client().StartSession(); err != nil {
		return err
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(sCtx mongo.SessionContext) (sErr error) {
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
	var (
		ctx     = service.ctx
		dbConn  = service.dbConn
		session mongo.Session
	)

	comment.ReplyCount--
	if session, err = dbConn.Client().StartSession(); err != nil {
		return err
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(sCtx mongo.SessionContext) (sErr error) {
		if err = repositories.DeleteComment(sCtx, dbConn, reply); err != nil {
			return nil
		}
		if err = repositories.UpdateComment(sCtx, dbConn, comment); err != nil {
			return err
		}

		return nil
	})
}
