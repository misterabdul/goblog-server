package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/misterabdul/goblog-server/internal/server"
)

// Run the main server.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	server := server.GetServer()
	server.Run(getAddress())
}

func getAddress() string {
	host := "localhost"
	if envHost, ok := os.LookupEnv("APP_HOST"); ok {
		host = envHost
	}
	port := "80"
	if envPort, ok := os.LookupEnv("APP_PORT"); ok {
		port = envPort
	}
	fmt.Println(host + ":" + port)
	return host + ":" + port
}
