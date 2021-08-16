package migration

import (
	"context"
	"log"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/pkg/utils"
)

func Fresh(ctx context.Context) {
	dbConn, err := database.GetDBConnDefault(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Client().Disconnect(ctx)

	utils.ConsolePrintlnYellow("Dropping database: " + dbConn.Name())
	if err := dbConn.Drop(ctx); err != nil {
		log.Fatal(err)
	}
	utils.ConsolePrintlnWhite("Dropped database : " + dbConn.Name())

	Migrate(ctx)
}
