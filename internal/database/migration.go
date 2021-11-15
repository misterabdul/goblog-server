package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/misterabdul/goblog-server/internal/database/migrations"
	"github.com/misterabdul/goblog-server/pkg/utils"
)

const migrationCollectionName = "migrations"

// Common migration functionalities.
type MigrationInterface interface {
	// Get unique migration name.
	Name() (collectionName string)
	// Do the migration.
	Up(ctx context.Context, dbConn *mongo.Database) (err error)
	// Undo the migration.
	Down(ctx context.Context, dbConn *mongo.Database) (err error)
}

// Model for migration data
type migrationModel struct {
	UID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Batch      int                `json:"batch"`
	Name       string             `json:"name"`
	MigratedAt primitive.DateTime `json:"migratedAt"`
}

// Get all registered migration
func getMigrations() (migrationInstances []MigrationInterface) {
	return []MigrationInterface{
		new(migrations.CreateUsersCollection),
		new(migrations.CreatePostsCollection),
		new(migrations.CreateCommentsCollection),
		new(migrations.CreateCategoryCollection),
		new(migrations.CreateNotificationCollection),
		new(migrations.CreateRevokedTokenCollection)}
}

// Run the migrations.
func Migrate(ctx context.Context, dbConn *mongo.Database) (err error) {
	var (
		migrationsData []*migrationModel
		migration      MigrationInterface
		batch          int
	)

	if migrationsData, batch, err = getMigrationsData(ctx, dbConn); err != nil {
		migrationsData = []*migrationModel{}
		batch = 1
	} else {
		batch++
	}
	for _, migration = range getMigrations() {
		if !alreadyMigrated(migration, migrationsData) {
			utils.ConsolePrintlnYellow("Migrating: " + migration.Name())
			if err = noteMigration(ctx, dbConn, batch, migration); err != nil {
				return err
			}
			if err = migration.Up(ctx, dbConn); err != nil {
				return err
			}
			utils.ConsolePrintlnWhite("Migrated : " + migration.Name())
		}
	}

	return nil
}

// Rollback the migrations.
func Rollback(ctx context.Context, dbConn *mongo.Database) (err error) {
	var (
		migrationsData []*migrationModel
		data           *migrationModel
		migration      MigrationInterface
		latestBatch    int
		isIn           bool
	)

	if migrationsData, latestBatch, err = getMigrationsData(ctx, dbConn); err != nil {
		migrationsData = []*migrationModel{}
	}
	for _, migration = range getMigrations() {
		if data, isIn = isInBatch(latestBatch, migration, migrationsData); isIn {
			utils.ConsolePrintlnYellow("Removing: " + migration.Name())
			if err = migration.Down(ctx, dbConn); err != nil {
				return err
			}
			if err = deleteMigrationNote(ctx, dbConn, data); err != nil {
				return err
			}
			utils.ConsolePrintlnWhite("Removed : " + migration.Name())
		}
	}

	return nil
}

// Create the migrations collection.
func createMigrationsCollection(ctx context.Context, dbConn *mongo.Database) (err error) {
	if err = dbConn.CreateCollection(ctx, migrationCollectionName); err != nil {
		return err
	}
	indexes := []mongo.IndexModel{{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}, {
		Keys:    bson.D{{Key: "createdAt", Value: -1}},
		Options: nil,
	}}
	if _, err = dbConn.Collection(migrationCollectionName).Indexes().
		CreateMany(ctx, indexes); err != nil {
		return err
	}

	return nil
}

// Get the migrations data from database.
func getMigrationsData(ctx context.Context, dbConn *mongo.Database) (
	migrationsData []*migrationModel,
	batch int,
	err error,
) {
	var (
		cursor *mongo.Cursor
		data   migrationModel
	)

	if cursor, err = dbConn.Collection(migrationCollectionName).
		Find(ctx, bson.M{}); err != nil {
		// Maybe there's no migrations collection, try to create one.
		if err = createMigrationsCollection(ctx, dbConn); err != nil {
			return nil, 1, err
		}

		return nil, 1, nil
	}
	batch = 1
	for cursor.Next(ctx) {
		if err = cursor.Decode(&data); err != nil {
			return nil, batch, err
		}
		if data.Batch > batch {
			batch = data.Batch
		}
		migrationsData = append(migrationsData, &data)
	}

	return migrationsData, batch, nil
}

// Check if given migration is already migrated.
func alreadyMigrated(migration MigrationInterface, data []*migrationModel) (alreadyMigrated bool) {
	for _, iData := range data {
		if migration.Name() == iData.Name {
			return true
		}
	}

	return false
}

// Save the migration data into database.
func noteMigration(
	ctx context.Context,
	dbConn *mongo.Database,
	batch int,
	migration MigrationInterface,
) (err error) {
	data := migrationModel{
		UID:        primitive.NewObjectID(),
		Batch:      batch,
		Name:       migration.Name(),
		MigratedAt: primitive.NewDateTimeFromTime(time.Now())}
	if _, err = dbConn.Collection(migrationCollectionName).
		InsertOne(ctx, data); err != nil {
		return err
	}

	return nil
}

// Check if given migration is in same batch of migration data.
func isInBatch(
	batch int,
	migration MigrationInterface,
	data []*migrationModel,
) (found *migrationModel, isInBatch bool) {
	for _, found = range data {
		if found.Batch != batch {
			continue
		}
		if migration.Name() == found.Name {
			return found, true
		}
	}

	return nil, false
}

// Remove migration data from database.
func deleteMigrationNote(
	ctx context.Context,
	dbConn *mongo.Database,
	data *migrationModel,
) (err error) {
	if _, err = dbConn.Collection(migrationCollectionName).
		DeleteOne(ctx, bson.M{"_id": data.UID}); err != nil {
		return err
	}

	return nil
}
