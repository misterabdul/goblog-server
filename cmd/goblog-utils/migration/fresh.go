package migration

import (
	"context"
	"log"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func Fresh(ctx context.Context) {
	var (
		dbConn *mongo.Database
		err    error
	)

	if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
		log.Fatal(err)
	}
	defer dbConn.Client().Disconnect(ctx)
	utils.ConsolePrintlnYellow("Dropping database: " + dbConn.Name())
	if err = dbConn.Drop(ctx); err != nil {
		log.Fatal(err)
	}
	utils.ConsolePrintlnWhite("Dropped database : " + dbConn.Name())
	Migrate(ctx)
}
