package migration

import (
	"context"
	"log"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func Rollback(ctx context.Context) {
	var (
		dbConn *mongo.Database
		err    error
	)

	if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
		log.Fatal(err)
	}
	defer dbConn.Client().Disconnect(ctx)
	if err = database.Rollback(ctx, dbConn); err != nil {
		log.Fatal(err)
	}
	utils.ConsolePrintlnGreen("Rollback migration completed.")
}
