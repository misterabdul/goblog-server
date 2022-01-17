package service

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

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
) (err error) {
	var (
		now = primitive.NewDateTimeFromTime(time.Now())
	)

	comment.UID = primitive.NewObjectID()
	comment.CreatedAt = now
	comment.DeletedAt = nil

	return repositories.SaveComment(
		service.ctx,
		service.dbConn,
		comment)
}

// Delete comment to trash
func (service *Service) TrashComment(
	comment *models.CommentModel,
) (err error) {
	now := primitive.NewDateTimeFromTime(time.Now())
	comment.DeletedAt = now

	return repositories.UpdateComment(
		service.ctx,
		service.dbConn,
		comment)
}

// Restore comment from trash
func (service *Service) DetrashComment(
	comment *models.CommentModel,
) (err error) {
	comment.DeletedAt = nil

	return repositories.UpdateComment(
		service.ctx,
		service.dbConn,
		comment)
}

// Permanently delete comment
func (service *Service) DeleteComment(
	comment *models.CommentModel,
) (err error) {

	return repositories.DeleteComment(
		service.ctx,
		service.dbConn,
		comment)
}
