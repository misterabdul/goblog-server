package migration

import (
	"context"
	"log"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/pkg/utils"
)

func Migrate(ctx context.Context) {
	dbConn, err := database.GetDBConnDefault(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Client().Disconnect(ctx)

	if err := database.Migrate(ctx, dbConn); err != nil {
		log.Fatal(err)
	}
	utils.ConsolePrintlnGreen("All migrations completed.")
}
