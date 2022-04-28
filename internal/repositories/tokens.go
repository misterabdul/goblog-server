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

type RevokedTokenRepository struct {
	collection *mongo.Collection
}

func NewRevokedTokenRepository(
	dbConn *mongo.Database,
) *RevokedTokenRepository {

	return &RevokedTokenRepository{
		collection: dbConn.Collection("revokedTokens")}
}

// Get single revoked tokens
func (r RevokedTokenRepository) ReadOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (revokedToken *models.RevokedTokenModel, err error) {
	var _revokedToken models.RevokedTokenModel

	if err = r.collection.FindOne(
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
func (r RevokedTokenRepository) ReadMany(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (revokedTokens []*models.RevokedTokenModel, err error) {
	var (
		revokedToken *models.RevokedTokenModel
		cursor       *mongo.Cursor
	)

	if cursor, err = r.collection.Find(
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
func (r RevokedTokenRepository) Save(
	ctx context.Context,
	revokedToken *models.RevokedTokenModel,
) (err error) {
	var (
		insRes     *mongo.InsertOneResult
		insertedID primitive.ObjectID
		ok         bool
	)

	if insRes, err = r.collection.InsertOne(
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
func (r RevokedTokenRepository) Update(
	ctx context.Context,
	revokedToken *models.RevokedTokenModel,
) (err error) {
	_, err = r.collection.UpdateByID(
		ctx, revokedToken.UID, bson.M{"$set": revokedToken})

	return err
}

// Delete revoked token
func (r RevokedTokenRepository) Delete(
	ctx context.Context,
	revokedToken *models.RevokedTokenModel,
) (err error) {
	_, err = r.collection.DeleteOne(
		ctx, bson.M{"_id": revokedToken.UID})

	return err
}
