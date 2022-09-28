package repositories

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/database/models"
)

const revokedTokenCollection = "revokedTokens"

// Get single revoked token
func ReadOneRevokedToken(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) (revokedToken *models.RevokedTokenModel, err error) {
	var (
		collection    = dbConn.Collection(revokedTokenCollection)
		_revokedToken models.RevokedTokenModel
	)

	if err = collection.FindOne(ctx, filter, opts...).Decode(&_revokedToken); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &_revokedToken, nil
}

// Get multiple revoked tokens
func ReadManyRevokedTokens(
	dbConn *mongo.Database,
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (revokedTokens []*models.RevokedTokenModel, err error) {
	var (
		collection   = dbConn.Collection(revokedTokenCollection)
		revokedToken *models.RevokedTokenModel
		cursor       *mongo.Cursor
	)

	if cursor, err = collection.Find(ctx, filter, opts...); err != nil {
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
func SaveOneRevokedToken(
	dbConn *mongo.Database,
	ctx context.Context,
	revokedToken *models.RevokedTokenModel,
	opts ...*options.InsertOneOptions,
) (err error) {
	var (
		collection = dbConn.Collection(revokedTokenCollection)
		insRes     *mongo.InsertOneResult
		insertedID primitive.ObjectID
		ok         bool
	)

	if insRes, err = collection.InsertOne(ctx, revokedToken, opts...); err != nil {
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
func UpdateOneRevokedToken(
	dbConn *mongo.Database,
	ctx context.Context,
	revokedToken *models.RevokedTokenModel,
	opts ...*options.UpdateOptions,
) (err error) {
	var collection = dbConn.Collection(revokedTokenCollection)

	_, err = collection.UpdateOne(
		ctx, bson.M{"_id": revokedToken.UID}, bson.M{"$set": revokedToken}, opts...)

	return err
}

// Delete revoked token
func DeleteOneRevokedToken(
	dbConn *mongo.Database,
	ctx context.Context,
	revokedToken *models.RevokedTokenModel,
	opts ...*options.DeleteOptions,
) (err error) {
	var collection = dbConn.Collection(revokedTokenCollection)

	_, err = collection.DeleteOne(
		ctx, bson.M{"_id": revokedToken.UID})

	return err
}
