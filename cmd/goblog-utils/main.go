package main

import (
	"bufio"
	"context"
	"os"
	"strings"

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
		utils.ConsolePrintlnWhite("No command specified")
		return
	}

	godotenv.Load()
	ctx := context.TODO()
	reader := bufio.NewReader(os.Stdin)

	switch {
	default:
		utils.ConsolePrintlnWhite("Unknown command : " + args)
	case args == "migration:fresh":
		utils.ConsolePrintlnYellow("Warning: this will create a new fresh database")
		utils.ConsolePrintWhite("Are you sure to execute migrate:fresh [yes/no]: ")
		if input, err := reader.ReadString('\n'); err == nil {
			if strings.ToLower(string(input)) == "yes\n" {
				migration.Fresh(ctx)
			}
		}
	case args == "migration:migrate":
		migration.Migrate(ctx)
	case args == "migration:rollback":
		migration.Rollback(ctx)
	}
}
