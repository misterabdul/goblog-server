package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"

	"github.com/misterabdul/goblog-server/cmd/goblog-utils/migration"
	"github.com/misterabdul/goblog-server/pkg/utils"
)

func main() {
	var args string
	if len(os.Args) > 1 {
		args = os.Args[1]
	}

	if len(args) < 1 {
		utils.ConsolePrintWhite("No command specified")
		return
	}

	godotenv.Load()

	ctx := context.TODO()

	switch {
	default:
		utils.ConsolePrintWhite("Unknown command : " + args)
	case args == "migration:migrate":
		migration.Migrate(ctx)
	case args == "migration:rollback":
		migration.Rollback(ctx)
	}
}
