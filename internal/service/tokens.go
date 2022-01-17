package service

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/models"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

// Get single revoked tokens
func (service *Service) GetRevokedToken(
	filter interface{},
) (revokedToken *models.RevokedTokenModel, err error) {
	return repositories.GetRevokedToken(
		service.ctx,
		service.dbConn,
		filter)
}

// Get multiple revoked tokens
func (service *Service) GetRevokedTokens(
	filter interface{},
) (revokedTokens []*models.RevokedTokenModel, err error) {
	return repositories.GetRevokedTokens(
		service.ctx,
		service.dbConn,
		filter,
		internalGin.GetFindOptions(service.c))
}

// Create new revoked token
func (service *Service) CreateRevokedToken(
	revokedToken *models.RevokedTokenModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	revokedToken.CreatedAt = now
	revokedToken.UpdatedAt = now
	revokedToken.DeletedAt = nil

	return repositories.SaveRevokedToken(
		service.ctx,
		service.dbConn,
		revokedToken)
}

// Update revoked token
func (service *Service) UpdateRevokedToken(
	revokedToken *models.RevokedTokenModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	revokedToken.UpdatedAt = now

	return repositories.UpdateRevokedToken(
		service.ctx,
		service.dbConn,
		revokedToken)
}

// Delete revoked to trash
func (service *Service) TrashRevokedToken(
	revokedToken *models.RevokedTokenModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	revokedToken.DeletedAt = now

	return repositories.UpdateRevokedToken(
		service.ctx,
		service.dbConn,
		revokedToken)
}

// Restore revoked token from trash
func (service *Service) DetrashRevokedToken(
	revokedToken *models.RevokedTokenModel,
) (err error) {
	revokedToken.DeletedAt = nil

	return repositories.UpdateRevokedToken(
		service.ctx,
		service.dbConn,
		revokedToken)
}

// Permanently delete revoked token
func (service *Service) DeleteRevokedToken(
	revokedToken *models.RevokedTokenModel,
) (err error) {

	return repositories.DeleteRevokedToken(
		service.ctx,
		service.dbConn,
		revokedToken)
}
