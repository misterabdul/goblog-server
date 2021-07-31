package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/pkg/utils"
)

func main() {
	godotenv.Load()

	ctx := context.TODO()

	dbConn, err := database.GetDBConnDefault(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Client().Disconnect(ctx)

	var args string
	if len(os.Args) > 1 {
		args = os.Args[1]
	}

	switch {
	default:
		fallthrough
	case args == "--migrate":
		if err := database.Migrate(ctx, dbConn); err != nil {
			log.Fatal(err)
		}
	case args == "--rollback":
		if err := database.Rollback(ctx, dbConn); err != nil {
			log.Fatal(err)
		}
	}

	utils.ConsolePrintGreen("All migrations completed.")
}
