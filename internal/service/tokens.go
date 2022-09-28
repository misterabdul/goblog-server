package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/database/repositories"
)

type revokedToken struct {
	dbConn *mongo.Database
}

func newRevokedTokenService(
	dbConn *mongo.Database,
) *revokedToken {

	return &revokedToken{dbConn: dbConn}
}

// Get single revoked tokens
func (s *revokedToken) GetOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (revokedToken *models.RevokedTokenModel, err error) {

	return repositories.ReadOneRevokedToken(
		s.dbConn, ctx, filter, opts...)
}

// Get multiple revoked tokens
func (s *revokedToken) GetMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (revokedTokens []*models.RevokedTokenModel, err error) {

	return repositories.ReadManyRevokedTokens(
		s.dbConn, ctx, filter, opts...)
}

// Create new revoked token
func (s *revokedToken) SaveOne(
	ctx context.Context,
	revokedToken *models.RevokedTokenModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	revokedToken.CreatedAt = now
	revokedToken.UpdatedAt = now
	revokedToken.DeletedAt = nil

	return repositories.SaveOneRevokedToken(
		s.dbConn, ctx, revokedToken, opts...)
}

// Update revoked token
func (s *revokedToken) UpdateOne(
	ctx context.Context,
	revokedToken *models.RevokedTokenModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	revokedToken.UpdatedAt = now

	return repositories.UpdateOneRevokedToken(
		s.dbConn, ctx, revokedToken, opts...)
}

// Delete revoked to trash
func (s *revokedToken) TrashRevokedToken(
	ctx context.Context,
	revokedToken *models.RevokedTokenModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	revokedToken.DeletedAt = now

	return repositories.UpdateOneRevokedToken(
		s.dbConn, ctx, revokedToken, opts...)
}

// Restore revoked token from trash
func (s *revokedToken) DetrashRevokedToken(
	ctx context.Context,
	revokedToken *models.RevokedTokenModel,
	opts ...*options.UpdateOptions,
) (err error) {
	revokedToken.DeletedAt = nil

	return repositories.UpdateOneRevokedToken(
		s.dbConn, ctx, revokedToken, opts...)
}

// Permanently delete revoked token
func (s *revokedToken) DeleteRevokedToken(
	ctx context.Context,
	revokedToken *models.RevokedTokenModel,
	opts ...*options.DeleteOptions,
) (err error) {

	return repositories.DeleteOneRevokedToken(
		s.dbConn, ctx, revokedToken, opts...)
}
