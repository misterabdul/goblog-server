package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database/migrations"
	"github.com/misterabdul/goblog-server/pkg/utils"
)

type MigrationInterface interface {
	Name() string
	Up(ctx context.Context, dbConn *mongo.Database) error
	Down(ctx context.Context, dbConn *mongo.Database) error
}

type migrationModel struct {
	UID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Batch      int                `json:"batch"`
	Name       string             `json:"name"`
	MigratedAt primitive.DateTime `json:"migratedAt"`
}

func getMigrations() []MigrationInterface {
	return []MigrationInterface{
		new(migrations.CreateUsersCollection),
	}
}

func getMigrationsData(ctx context.Context, dbConn *mongo.Database) ([]*migrationModel, int, error) {
	cursor, err := dbConn.Collection("migrations").Find(ctx, bson.M{})
	if err != nil {
		return nil, 1, err
	}

	var migrationsData []*migrationModel
	batch := 1
	for cursor.Next(ctx) {
		var data migrationModel
		if err := cursor.Decode(&data); err != nil {
			return nil, batch, err
		}
		if data.Batch > batch {
			batch = data.Batch
		}
		migrationsData = append(migrationsData, &data)
	}

	return migrationsData, batch, nil
}

func alreadyMigrated(migration MigrationInterface, data []*migrationModel) bool {
	for _, iData := range data {
		if migration.Name() == iData.Name {
			return true
		}
	}
	return false
}

func noteMigration(ctx context.Context, dbConn *mongo.Database, batch int, migration MigrationInterface) error {
	data := migrationModel{
		UID:        primitive.NewObjectID(),
		Batch:      batch,
		Name:       migration.Name(),
		MigratedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
	_, err := dbConn.Collection("migrations").InsertOne(ctx, data)

	return err
}

func Migrate(ctx context.Context, dbConn *mongo.Database) error {
	migrationsData, batch, err := getMigrationsData(ctx, dbConn)
	if err != nil {
		migrationsData = []*migrationModel{}
		batch = 1
	} else {
		batch++
	}

	for _, migration := range getMigrations() {
		if !alreadyMigrated(migration, migrationsData) {
			utils.ConsolePrintYellow("Migrating: " + migration.Name())
			if err := noteMigration(ctx, dbConn, batch, migration); err != nil {
				return err
			}
			if err := migration.Up(ctx, dbConn); err != nil {
				return err
			}
			utils.ConsolePrintGreen("Migrated: " + migration.Name())
		}
	}

	return nil
}

func isInBatch(batch int, migration MigrationInterface, data []*migrationModel) (*migrationModel, bool) {
	for _, iData := range data {
		if iData.Batch != batch {
			continue
		}
		if migration.Name() == iData.Name {
			return iData, true
		}
	}
	return nil, false
}

func deleteMigrationNote(ctx context.Context, dbConn *mongo.Database, data *migrationModel) error {
	_, err := dbConn.Collection("migrations").DeleteOne(ctx, bson.M{"_id": data.UID})

	return err
}

func Rollback(ctx context.Context, dbConn *mongo.Database) error {
	migrationsData, latestBatch, err := getMigrationsData(ctx, dbConn)
	if err != nil {
		migrationsData = []*migrationModel{}
	}

	for _, migration := range getMigrations() {
		if data, isIn := isInBatch(latestBatch, migration, migrationsData); isIn {
			utils.ConsolePrintYellow("Removing: " + migration.Name())
			if err := migration.Down(ctx, dbConn); err != nil {
				return err
			}
			if err := deleteMigrationNote(ctx, dbConn, data); err != nil {
				return err
			}
			utils.ConsolePrintGreen("Removed: " + migration.Name())
		}
	}

	return nil
}
