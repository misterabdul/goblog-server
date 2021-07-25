package main

import (
	"os"

	"github.com/misterabdul/goblog-server/internal/server"
)

// Run the main server.
func main() {
	server := server.GetServer()
	server.Run(os.Args[1])
}
