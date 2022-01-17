package repositories

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/models"
)

func getRevokedTokenCollection(
	dbConn *mongo.Database,
) (tokenCollection *mongo.Collection) {
	return dbConn.Collection("revokedTokens")
}

// Get single revoked tokens
func GetRevokedToken(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.FindOneOptions,
) (revokedToken *models.RevokedTokenModel, err error) {
	var _revokedToken models.RevokedTokenModel

	if err = getRevokedTokenCollection(dbConn).FindOne(
		ctx, filter, opts...,
	).Decode(&_revokedToken); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_revokedToken, nil
}

// Get multiple revoked tokens
func GetRevokedTokens(
	ctx context.Context,
	dbConn *mongo.Database,
	filter interface{},
	opts ...*options.FindOptions,
) (revokedTokens []*models.RevokedTokenModel, err error) {
	var (
		revokedToken *models.RevokedTokenModel
		cursor       *mongo.Cursor
	)

	if cursor, err = getRevokedTokenCollection(dbConn).Find(
		ctx, filter, opts...,
	); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		revokedToken = &models.RevokedTokenModel{}
		if err := cursor.Decode(revokedToken); err != nil {
			return nil, err
		}
		revokedTokens = append(revokedTokens, revokedToken)
	}

	return revokedTokens, nil
}

// Save new revoked token
func SaveRevokedToken(
	ctx context.Context,
	dbConn *mongo.Database,
	revokedToken *models.RevokedTokenModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID primitive.ObjectID
		ok         bool
	)

	if insRes, err = getRevokedTokenCollection(dbConn).InsertOne(
		ctx, revokedToken,
	); err != nil {
		return err
	}
	if insertedID, ok = insRes.InsertedID.(primitive.ObjectID); !ok {
		return errors.New("unable to assert inserted uid")
	}
	if revokedToken.UID != insertedID {
		return errors.New("inserted uid is not same with database")
	}

	return nil
}

// Update revoked token
func UpdateRevokedToken(
	ctx context.Context,
	dbConn *mongo.Database,
	revokedToken *models.RevokedTokenModel,
) (err error) {
	_, err = getRevokedTokenCollection(dbConn).UpdateByID(
		ctx, revokedToken.UID, bson.M{"$set": revokedToken})

	return err
}

// Delete revoked token
func DeleteRevokedToken(
	ctx context.Context,
	dbConn *mongo.Database,
	revokedToken *models.RevokedTokenModel,
) (err error) {
	_, err = getRevokedTokenCollection(dbConn).DeleteOne(
		ctx, bson.M{"_id": revokedToken.UID})

	return err
}
