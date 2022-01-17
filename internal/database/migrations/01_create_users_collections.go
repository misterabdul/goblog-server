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

const usersCollectionName = "users"

// Create the users collection.
// Insert super admin record data.
type CreateUsersCollection struct{}

func (m *CreateUsersCollection) Name() (collectionName string) {
	return "01_create_users_collections"
}

func (m *CreateUsersCollection) Up(ctx context.Context, dbConn *mongo.Database) (err error) {
	if err = dbConn.CreateCollection(ctx, usersCollectionName); err != nil {
		return err
	}
	indexes := []mongo.IndexModel{{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}, {
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}, {
		Keys:    bson.D{{Key: "createdAt", Value: -1}},
		Options: nil,
	}, {
		Keys:    bson.D{{Key: "updatedAt", Value: -1}},
		Options: nil,
	}, {
		Keys:    bson.D{{Key: "deletedAt", Value: -1}},
		Options: nil,
	}}
	if _, err = dbConn.Collection(usersCollectionName).Indexes().
		CreateMany(ctx, indexes); err != nil {
		return err
	}
	insertSuperAdmin(ctx, dbConn)

	return nil
}

func (m *CreateUsersCollection) Down(ctx context.Context, dbConn *mongo.Database) (err error) {
	return dbConn.Collection(usersCollectionName).Drop(ctx)
}

func insertSuperAdmin(ctx context.Context, dbConn *mongo.Database) (err error) {
	var (
		password   string
		now        = primitive.NewDateTimeFromTime(time.Now())
		superAdmin models.UserModel
	)

	if password, err = hash.Make("password"); err != nil {
		return err
	}
	superAdmin = models.UserModel{
		FirstName: "Super Admin",
		Email:     "superadmin@example.com",
		Username:  "superadmin",
		Password:  password,
		Roles: []models.UserRole{{
			Level: 0,
			Name:  "SuperAdmin",
			Since: now,
		}, {
			Level: 1,
			Name:  "Admin",
			Since: now,
		}, {
			Level: 2,
			Name:  "Editor",
			Since: now,
		}, {
			Level: 3,
			Name:  "Writer",
			Since: now,
		}},
		CreatedAt: now,
		UpdatedAt: now,
	}

	return repositories.SaveUser(ctx, dbConn, &superAdmin)
}
