package main

import (
	"os"

	"github.com/misterabdul/goblog-server/internal/server"
)

// Run the main server.
func main() {
	server := server.GetServer()
	args := ""
	if len(os.Args) > 1 {
		args = os.Args[1]
	}
	server.Run(args)
}
