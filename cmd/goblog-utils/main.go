package main

import (
	"bufio"
	"context"
	"os"
	"sort"
	"strings"

	"github.com/joho/godotenv"

	"github.com/misterabdul/goblog-server/cmd/goblog-utils/fakedata"
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

	availableCommands := getAvailableCommands()

	if commands, ok := availableCommands[args]; ok {
		commands(ctx, reader)
	} else {
		utils.ConsolePrintlnWhite("Unknown command: " + args)
		utils.ConsolePrintlnWhite("Available commands")
		availableCommandNames := getAvailableCommandNames(availableCommands)
		for _, name := range availableCommandNames {
			utils.ConsolePrintlnWhite("  " + name)
		}
	}

}

func getAvailableCommands() (
	availableCommands map[string]func(
		context.Context, *bufio.Reader),
) {
	return map[string]func(context.Context, *bufio.Reader){
		"migrations:fresh": func(ctx context.Context, reader *bufio.Reader) {
			utils.ConsolePrintlnYellow("Warning: this will create a new fresh database")
			utils.ConsolePrintWhite("Are you sure to execute migrate:fresh (y/N): ")
			if input, err := reader.ReadString('\n'); err == nil {
				if strings.ToLower(string(input)) == "y\n" {
					migration.Fresh(ctx)
				}
			}
		},
		"migrations:migrate": func(ctx context.Context, reader *bufio.Reader) {
			migration.Migrate(ctx)
		},
		"migrations:rollback": func(ctx context.Context, reader *bufio.Reader) {
			migration.Rollback(ctx)
		},
		"migrations:status": func(ctx context.Context, reader *bufio.Reader) {
			migration.Status(ctx)
		},
		"fakedata:user": func(ctx context.Context, reader *bufio.Reader) {
			fakedata.GenerateUsers(ctx)
		},
		"fakedata:category": func(ctx context.Context, reader *bufio.Reader) {
			fakedata.GenerateCategories(ctx)
		},
		"fakedata:post": func(ctx context.Context, reader *bufio.Reader) {
			fakedata.GeneratePosts(ctx)
		},
		"fakedata:page": func(ctx context.Context, reader *bufio.Reader) {
			fakedata.GeneratePages(ctx)
		},
	}
}

func getAvailableCommandNames(
	availableCommands map[string]func(context.Context, *bufio.Reader),
) []string {
	var names []string

	names = make([]string, 0, len(availableCommands))
	for keys := range availableCommands {
		names = append(names, keys)
	}
	sort.Strings(names)

	return names
}
