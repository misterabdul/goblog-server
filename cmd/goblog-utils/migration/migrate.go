package migration

import (
	"context"
	"log"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func Migrate(ctx context.Context) {
	var (
		dbConn *mongo.Database
		err    error
	)

	if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
		log.Fatal(err)
	}
	defer dbConn.Client().Disconnect(ctx)
	if err = database.Migrate(ctx, dbConn); err != nil {
		log.Fatal(err)
	}
	utils.ConsolePrintlnGreen("All migrations completed.")
}
