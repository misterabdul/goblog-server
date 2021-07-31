package migrations

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/pkg/hash"
)

// Create the users collection.
// Insert super admin record data.
type CreateUsersCollection struct {
}

func (m *CreateUsersCollection) Name() string {
	return "01_create_users_collections"
}

func (m *CreateUsersCollection) Up(ctx context.Context, dbConn *mongo.Database) error {
	if err := dbConn.CreateCollection(ctx, "users"); err != nil {
		return err
	}

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "createdAt", Value: -1}},
			Options: nil,
		},
	}

	_, err := dbConn.Collection("users").Indexes().CreateMany(ctx, indexes)

	insertSuperAdmin(ctx, dbConn)

	return err
}

func (m *CreateUsersCollection) Down(ctx context.Context, dbConn *mongo.Database) error {
	return dbConn.Collection("users").Drop(ctx)
}

func insertSuperAdmin(ctx context.Context, dbConn *mongo.Database) error {
	password, err := hash.Make("password")
	if err != nil {
		return err
	}

	now := time.Now()
	superAdmin := models.UserModel{
		UID:       primitive.NewObjectID(),
		FirstName: "Super Admin",
		Email:     "superadmin@example.com",
		Username:  "superadmin",
		Password:  password,
		CreatedAt: primitive.NewDateTimeFromTime(now),
		UpdatedAt: primitive.NewDateTimeFromTime(now),
	}

	return repositories.CreateUser(ctx, dbConn, &superAdmin)
}