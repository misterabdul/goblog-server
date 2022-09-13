package service

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/database/repositories"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
)

type RevokedTokenService struct {
	c          *gin.Context
	ctx        context.Context
	dbConn     *mongo.Database
	repository *repositories.RevokedTokenRepository
}

func NewRevokedTokenService(
	c *gin.Context,
	ctx context.Context,
	dbConn *mongo.Database,
) *RevokedTokenService {

	return &RevokedTokenService{
		c:          c,
		ctx:        ctx,
		dbConn:     dbConn,
		repository: repositories.NewRevokedTokenRepository(dbConn)}
}

// Get single revoked tokens
func (s *RevokedTokenService) GetRevokedToken(
	filter interface{},
) (revokedToken *models.RevokedTokenModel, err error) {
	return s.repository.ReadOne(
		s.ctx, filter)
}

// Get multiple revoked tokens
func (s *RevokedTokenService) GetRevokedTokens(
	filter interface{},
) (revokedTokens []*models.RevokedTokenModel, err error) {
	return s.repository.ReadMany(
		s.ctx, filter,
		internalGin.GetFindOptions(s.c))
}

// Create new revoked token
func (s *RevokedTokenService) CreateRevokedToken(
	revokedToken *models.RevokedTokenModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	revokedToken.CreatedAt = now
	revokedToken.UpdatedAt = now
	revokedToken.DeletedAt = nil

	return s.repository.Save(
		s.ctx, revokedToken)
}

// Update revoked token
func (s *RevokedTokenService) UpdateRevokedToken(
	revokedToken *models.RevokedTokenModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	revokedToken.UpdatedAt = now

	return s.repository.Update(
		s.ctx, revokedToken)
}

// Delete revoked to trash
func (s *RevokedTokenService) TrashRevokedToken(
	revokedToken *models.RevokedTokenModel,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	revokedToken.DeletedAt = now

	return s.repository.Update(
		s.ctx, revokedToken)
}

// Restore revoked token from trash
func (s *RevokedTokenService) DetrashRevokedToken(
	revokedToken *models.RevokedTokenModel,
) (err error) {
	revokedToken.DeletedAt = nil

	return s.repository.Update(
		s.ctx, revokedToken)
}

// Permanently delete revoked token
func (s *RevokedTokenService) DeleteRevokedToken(
	revokedToken *models.RevokedTokenModel,
) (err error) {

	return s.repository.Delete(
		s.ctx, revokedToken)
}
